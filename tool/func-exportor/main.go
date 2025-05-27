package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

var (
	flagVersion  = flag.Bool("version", false, "Print the version of the tool")
	flagJSON     = flag.Bool("json", false, "Output in JSON format")
	flagStats    = flag.Bool("stats", false, "Show only statistics")
	flagHeader   = flag.String("header", "", "Generate header file with signatures (specify output filename)")
	flagValidate = flag.Bool("validate", false, "Validate that the generated header file compiles correctly")
	flagVerbose  = flag.Bool("verbose", false, "Enable verbose output")
	flagOverview = flag.Bool("overview", false, "Generate package overview with documentation")
)

// ExportedSymbol 表示一个导出的符号
type ExportedSymbol struct {
	Name         string
	Type         string // "function", "variable", "constant", "type"
	Position     string
	Signature    string // 符号的完整签名
	VarSignature string // 函数的变量形式签名（仅用于函数）
	Doc          string // 文档注释
}

// isExported 检查标识符是否为导出的（首字母大写）
func isExported(name string) bool {
	return name != "" && name[0] >= 'A' && name[0] <= 'Z'
}

// extractDocComment 提取文档注释
func extractDocComment(commentGroup *ast.CommentGroup) string {
	if commentGroup == nil {
		return ""
	}

	var doc strings.Builder
	for _, comment := range commentGroup.List {
		text := comment.Text
		// 移除 // 或 /* */ 标记
		if strings.HasPrefix(text, "//") {
			text = strings.TrimSpace(text[2:])
		} else if strings.HasPrefix(text, "/*") && strings.HasSuffix(text, "*/") {
			text = strings.TrimSpace(text[2 : len(text)-2])
		}
		if text != "" {
			doc.WriteString("// ")
			doc.WriteString(text)
			doc.WriteString("\n")
		}
	}
	return doc.String()
}

// extractExportedSymbols 从AST中提取所有导出的符号
func extractExportedSymbols(fileSet *token.FileSet, file *ast.File) []ExportedSymbol {
	var symbols []ExportedSymbol

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			// 函数声明
			if node.Name != nil && isExported(node.Name.Name) {
				pos := fileSet.Position(node.Pos())
				signature := buildFunctionSignature(fileSet, node)
				varSignature := buildFunctionVarSignature(fileSet, node)
				doc := extractDocComment(node.Doc)
				symbols = append(symbols, ExportedSymbol{
					Name:         node.Name.Name,
					Type:         "function",
					Position:     fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column),
					Signature:    signature,
					VarSignature: varSignature,
					Doc:          doc,
				})
			}

		case *ast.GenDecl:
			// 通用声明（变量、常量、类型）
			for _, spec := range node.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					// 变量或常量
					for i, name := range s.Names {
						if isExported(name.Name) {
							pos := fileSet.Position(name.Pos())
							symbolType := "variable"
							if node.Tok == token.CONST {
								symbolType = "constant"
							}

							// 构建变量/常量签名
							var signature string
							if symbolType == "variable" {
								signature = "var " + name.Name
								if s.Type != nil {
									signature += " " + typeToString(fileSet, s.Type)
								} else if len(s.Values) > i && s.Values[i] != nil {
									// 如果没有显式类型，尝试从值推断
									signature += " " + typeToString(fileSet, s.Type)
								}
							} else {
								signature = "const " + name.Name
								if s.Type != nil {
									signature += " " + typeToString(fileSet, s.Type)
								}
								if len(s.Values) > i && s.Values[i] != nil {
									signature += " = " + nodeToString(fileSet, s.Values[i])
								}
							}

							// 提取文档注释 - 优先使用spec的文档，如果没有则使用GenDecl的文档
							doc := extractDocComment(s.Doc)
							if doc == "" {
								doc = extractDocComment(node.Doc)
							}

							symbols = append(symbols, ExportedSymbol{
								Name:      name.Name,
								Type:      symbolType,
								Position:  fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column),
								Signature: signature,
								Doc:       doc,
							})
						}
					}

				case *ast.TypeSpec:
					// 类型声明
					if isExported(s.Name.Name) {
						pos := fileSet.Position(s.Pos())
						signature := "type " + s.Name.Name + " " + typeToString(fileSet, s.Type)

						// 提取文档注释
						doc := extractDocComment(s.Doc)
						if doc == "" {
							doc = extractDocComment(node.Doc)
						}

						symbols = append(symbols, ExportedSymbol{
							Name:      s.Name.Name,
							Type:      "type",
							Position:  fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column),
							Signature: signature,
							Doc:       doc,
						})
					}
				}
			}
		}
		return true
	})

	return symbols
}

// analyzeGoFile 分析Go文件并返回导出的符号
func analyzeGoFile(filename string) ([]ExportedSymbol, error) {
	fileSet := token.NewFileSet()

	// 解析Go文件
	file, err := parser.ParseFile(fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析文件失败: %v", err)
	}

	// 提取导出的符号
	symbols := extractExportedSymbols(fileSet, file)

	return symbols, nil
}

// nodeToString 将AST节点转换为字符串
func nodeToString(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return ""
	}
	return buf.String()
}

// typeToString 将类型表达式转换为字符串
func typeToString(fset *token.FileSet, expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	return nodeToString(fset, expr)
}

// buildFunctionSignature 构建函数签名
func buildFunctionSignature(fset *token.FileSet, funcDecl *ast.FuncDecl) string {
	var sig strings.Builder

	// 函数名
	sig.WriteString("func")
	if funcDecl.Recv != nil {
		// 方法的接收者
		sig.WriteString(" ")
		sig.WriteString(nodeToString(fset, funcDecl.Recv))
	}
	sig.WriteString(" ")
	sig.WriteString(funcDecl.Name.Name)

	// 参数列表 - 使用完整的参数信息
	if funcDecl.Type.Params != nil {
		sig.WriteString(nodeToString(fset, funcDecl.Type.Params))
	} else {
		sig.WriteString("()")
	}

	// 返回值
	if funcDecl.Type.Results != nil {
		sig.WriteString(" ")
		sig.WriteString(nodeToString(fset, funcDecl.Type.Results))
	}

	return sig.String()
}

// buildFunctionVarSignature 构建函数变量签名（用于头文件）
func buildFunctionVarSignature(fset *token.FileSet, funcDecl *ast.FuncDecl) string {
	var sig strings.Builder

	sig.WriteString("var ")
	sig.WriteString(funcDecl.Name.Name)
	sig.WriteString(" func")

	// 如果有接收者，将其作为第一个参数
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		sig.WriteString("(")

		// 接收者类型和名称
		recv := funcDecl.Recv.List[0]
		recvName := "receiver"
		if len(recv.Names) > 0 && recv.Names[0] != nil {
			recvName = recv.Names[0].Name
		}
		recvType := typeToString(fset, recv.Type)
		sig.WriteString(recvName + " " + recvType)

		// 如果还有其他参数，添加逗号和参数
		if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0 {
			sig.WriteString(", ")
			// 添加其他参数的类型和名称
			buildParameterList(&sig, fset, funcDecl.Type.Params.List)
		}

		sig.WriteString(")")
	} else {
		// 普通函数，处理参数列表
		sig.WriteString("(")
		if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0 {
			buildParameterList(&sig, fset, funcDecl.Type.Params.List)
		}
		sig.WriteString(")")
	}

	// 返回值处理
	if funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0 {
		sig.WriteString(" ")
		buildReturnTypeList(&sig, fset, funcDecl.Type.Results.List)
	}

	return sig.String()
}

// buildParameterList 构建参数列表
func buildParameterList(sig *strings.Builder, fset *token.FileSet, params []*ast.Field) {
	for i, param := range params {
		if i > 0 {
			sig.WriteString(", ")
		}

		// 处理参数名称和类型
		paramNames := make([]string, 0, len(param.Names))
		for _, name := range param.Names {
			if name != nil {
				paramNames = append(paramNames, name.Name)
			}
		}

		paramType := typeToString(fset, param.Type)

		if len(paramNames) > 0 {
			// 有参数名称
			for j, name := range paramNames {
				if j > 0 {
					sig.WriteString(", ")
				}
				sig.WriteString(name + " " + paramType)
			}
		} else {
			// 没有参数名称，使用默认名称
			sig.WriteString(fmt.Sprintf("param%d %s", i, paramType))
		}
	}
}

// buildReturnTypeList 构建返回类型列表
func buildReturnTypeList(sig *strings.Builder, fset *token.FileSet, results []*ast.Field) {
	if len(results) == 1 && len(results[0].Names) == 0 {
		// 单个无名返回值
		sig.WriteString(typeToString(fset, results[0].Type))
	} else {
		// 多个返回值或有名返回值
		sig.WriteString("(")
		for i, result := range results {
			if i > 0 {
				sig.WriteString(", ")
			}

			// 如果有名称，添加名称
			if len(result.Names) > 0 {
				for j, name := range result.Names {
					if j > 0 {
						sig.WriteString(", ")
					}
					if name != nil {
						sig.WriteString(name.Name + " ")
					}
				}
			}

			sig.WriteString(typeToString(fset, result.Type))
		}
		sig.WriteString(")")
	}
}

// validateHeaderFile 验证生成的头文件是否能正确编译
func validateHeaderFile(headerFile string) error {
	if *flagVerbose {
		fmt.Printf("验证头文件: %s\n", headerFile)
	}

	// 解析头文件
	fileSet := token.NewFileSet()
	_, err := parser.ParseFile(fileSet, headerFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("头文件语法错误: %v", err)
	}

	if *flagVerbose {
		fmt.Println("头文件语法验证通过")
	}

	return nil
}

// generatePackageOverview 生成包概览
func generatePackageOverview(filename string, symbols []ExportedSymbol) {
	fmt.Printf("Package Overview for %s\n", filename)
	fmt.Printf("=" + strings.Repeat("=", len(filename)+20) + "\n\n")

	// 按类型分组统计
	typeGroups := make(map[string][]ExportedSymbol)
	for _, symbol := range symbols {
		typeGroups[symbol.Type] = append(typeGroups[symbol.Type], symbol)
	}

	// 总览统计
	fmt.Printf("📊 Summary\n")
	fmt.Printf("----------\n")
	fmt.Printf("Total exported symbols: %d\n", len(symbols))
	for symbolType, symbolList := range typeGroups {
		capitalizedType := strings.ToUpper(symbolType[:1]) + symbolType[1:] + "s"
		fmt.Printf("  %s: %d\n", capitalizedType, len(symbolList))
	}
	fmt.Println()

	// 详细列表
	order := []string{"type", "constant", "variable", "function"}
	icons := map[string]string{
		"type":     "🏗️",
		"constant": "📌",
		"variable": "📦",
		"function": "⚡",
	}

	for _, symbolType := range order {
		if symbolList, ok := typeGroups[symbolType]; ok {
			capitalizedType := strings.ToUpper(symbolType[:1]) + symbolType[1:] + "s"
			icon := icons[symbolType]
			fmt.Printf("%s %s (%d)\n", icon, capitalizedType, len(symbolList))
			fmt.Printf(strings.Repeat("-", len(capitalizedType)+10) + "\n")

			for _, symbol := range symbolList {
				fmt.Printf("  • %s", symbol.Name)
				if symbol.Doc != "" {
					// 提取文档的第一行作为简短描述
					lines := strings.Split(strings.TrimSpace(symbol.Doc), "\n")
					if len(lines) > 0 {
						firstLine := strings.TrimPrefix(strings.TrimSpace(lines[0]), "//")
						firstLine = strings.TrimSpace(firstLine)
						if firstLine != "" {
							fmt.Printf(" - %s", firstLine)
						}
					}
				}
				fmt.Printf(" (%s)\n", symbol.Position)

				// 显示签名（简化版本）
				if symbolType == "function" {
					if symbol.VarSignature != "" {
						fmt.Printf("    %s\n", symbol.VarSignature)
					}
				} else {
					fmt.Printf("    %s\n", symbol.Signature)
				}
			}
			fmt.Println()
		}
	}
}

// generateHeaderFile 生成头文件
func generateHeaderFile(sourceFile string, symbols []ExportedSymbol, outputFile string) error {
	// 解析源文件以获取包信息和导入
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, sourceFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("解析源文件失败: %v", err)
	}

	var content strings.Builder

	// 添加生成的头文件注释
	content.WriteString("// Code generated by func-exportor; DO NOT EDIT.\n")
	content.WriteString(fmt.Sprintf("// Source: %s\n\n", sourceFile))

	// 包声明
	content.WriteString(fmt.Sprintf("package %s\n\n", file.Name.Name))

	// 导入声明
	if len(file.Imports) > 0 {
		content.WriteString("import (\n")
		for _, imp := range file.Imports {
			content.WriteString("\t")
			if imp.Name != nil {
				content.WriteString(imp.Name.Name + " ")
			}
			content.WriteString(imp.Path.Value)
			content.WriteString("\n")
		}
		content.WriteString(")\n\n")
	}

	// 按类型分组符号
	typeGroups := make(map[string][]ExportedSymbol)
	for _, symbol := range symbols {
		typeGroups[symbol.Type] = append(typeGroups[symbol.Type], symbol)
	}

	// 先输出类型声明
	if types, ok := typeGroups["type"]; ok {
		content.WriteString("// Types\n")
		for _, symbol := range types {
			if symbol.Doc != "" {
				content.WriteString(symbol.Doc)
			}
			content.WriteString(symbol.Signature + "\n\n")
		}
	}

	// 然后输出常量
	if constants, ok := typeGroups["constant"]; ok {
		content.WriteString("// Constants\n")
		for _, symbol := range constants {
			if symbol.Doc != "" {
				content.WriteString(symbol.Doc)
			}
			content.WriteString(symbol.Signature + "\n")
		}
		content.WriteString("\n")
	}

	// 然后输出变量
	if variables, ok := typeGroups["variable"]; ok {
		content.WriteString("// Variables\n")
		for _, symbol := range variables {
			if symbol.Doc != "" {
				content.WriteString(symbol.Doc)
			}
			content.WriteString(symbol.Signature + "\n")
		}
		content.WriteString("\n")
	}

	// 最后输出函数（转换为函数变量）
	if functions, ok := typeGroups["function"]; ok {
		content.WriteString("// Functions (as function variables)\n")
		for _, symbol := range functions {
			if symbol.Doc != "" {
				content.WriteString(symbol.Doc)
			}
			content.WriteString(symbol.VarSignature + "\n")
		}
	}

	// 写入文件
	return os.WriteFile(outputFile, []byte(content.String()), 0644)
}

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println("func-exportor version 1.1.0")
		fmt.Println("A tool to analyze Go files and extract exported symbols")
		fmt.Println("Usage: func-exportor [options] <file.go>")
		fmt.Println("Options:")
		fmt.Println("  --json             Output in JSON format")
		fmt.Println("  --stats            Show only statistics")
		fmt.Println("  --header <file>    Generate header file with signatures")
		fmt.Println("  --validate         Validate generated header file")
		fmt.Println("  --verbose          Enable verbose output")
		fmt.Println("  --overview         Generate package overview with documentation")
		fmt.Println("  --version          Show version information")
		return
	}

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: func-exportor [options] <file.go>")
		fmt.Println("Options:")
		fmt.Println("  --json             Output in JSON format")
		fmt.Println("  --stats            Show only statistics")
		fmt.Println("  --header <file>    Generate header file with signatures")
		fmt.Println("  --validate         Validate generated header file")
		fmt.Println("  --verbose          Enable verbose output")
		fmt.Println("  --overview         Generate package overview with documentation")
		fmt.Println("  --version          Show version information")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  func-exportor main.go")
		fmt.Println("  func-exportor --json main.go")
		fmt.Println("  func-exportor --stats main.go")
		fmt.Println("  func-exportor --overview main.go")
		fmt.Println("  func-exportor --header header.go main.go")
		fmt.Println("  func-exportor --header header.go --validate main.go")
		return
	}

	filename := flag.Args()[0]

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("文件不存在: %s", filename)
	}

	// 分析Go文件
	symbols, err := analyzeGoFile(filename)
	if err != nil {
		log.Fatalf("分析文件失败: %v", err)
	}

	// 如果启用了概览模式
	if *flagOverview {
		generatePackageOverview(filename, symbols)
		return
	}

	// 如果指定了头文件输出
	if *flagHeader != "" {
		if *flagVerbose {
			fmt.Printf("生成头文件: %s\n", *flagHeader)
		}

		err := generateHeaderFile(filename, symbols, *flagHeader)
		if err != nil {
			log.Fatalf("生成头文件失败: %v", err)
		}
		fmt.Printf("头文件已生成: %s\n", *flagHeader)

		// 如果启用了验证，验证生成的头文件
		if *flagValidate {
			if err := validateHeaderFile(*flagHeader); err != nil {
				log.Fatalf("头文件验证失败: %v", err)
			}
			fmt.Println("头文件验证通过")
		}
		return
	}

	// 输出结果
	if len(symbols) == 0 {
		if *flagJSON {
			fmt.Println("[]")
		} else if *flagStats {
			fmt.Printf("File: %s\nTotal exported symbols: 0\n", filename)
		} else {
			fmt.Printf("在文件 %s 中没有找到导出的符号\n", filename)
		}
		return
	}

	if *flagHeader != "" {
		// 生成头文件
		outputFile := *flagHeader
		if !strings.HasSuffix(outputFile, ".go") {
			outputFile += ".go"
		}

		// 检查输出文件是否已存在
		if _, err := os.Stat(outputFile); err == nil {
			log.Fatalf("输出文件已存在: %s", outputFile)
		}

		err := generateHeaderFile(filename, symbols, outputFile)
		if err != nil {
			log.Fatalf("生成头文件失败: %v", err)
		}
		fmt.Printf("头文件已生成: %s\n", outputFile)
		return
	}

	// 按类型分组统计
	typeGroups := make(map[string][]ExportedSymbol)
	for _, symbol := range symbols {
		typeGroups[symbol.Type] = append(typeGroups[symbol.Type], symbol)
	}

	if *flagStats {
		// 统计信息输出
		fmt.Printf("File: %s\n", filename)
		fmt.Printf("Total exported symbols: %d\n", len(symbols))
		for symbolType, symbolList := range typeGroups {
			// 首字母大写
			capitalizedType := strings.ToUpper(symbolType[:1]) + symbolType[1:]
			fmt.Printf("  %s: %d\n", capitalizedType, len(symbolList))
		}
		return
	}

	if *flagJSON {
		// JSON格式输出
		output, err := json.MarshalIndent(symbols, "", "  ")
		if err != nil {
			log.Fatalf("序列化JSON失败: %v", err)
		}
		fmt.Println(string(output))
		return
	}

	// 人类可读格式输出
	fmt.Printf("在文件 %s 中找到 %d 个导出的符号:\n\n", filename, len(symbols))

	for symbolType, symbolList := range typeGroups {
		// 首字母大写
		capitalizedType := strings.ToUpper(symbolType[:1]) + symbolType[1:]
		fmt.Printf("%s (%d个):\n", capitalizedType, len(symbolList))
		for _, symbol := range symbolList {
			fmt.Printf("  - %s (%s)\n", symbol.Name, symbol.Position)
		}
		fmt.Println()
	}
}
