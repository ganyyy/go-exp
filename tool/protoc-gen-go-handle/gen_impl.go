package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"google.golang.org/protobuf/compiler/protogen"
)

func (s *Service) GenImpl(gen *protogen.Plugin) {
	path := filepath.Join(moduleBasePath, s.Package, s.Name+".impl.go")
	out := s.OutputFile(gen, ".impl")

	// 如果文件存在
	if fileIsExist(path) {
		// AST 增量
		s.GenFromAst(path, out)
	} else {
		// 当前不存在, 全量生成
		out.P("package ", s.Package)
		out.P()
		out.P("var (")
		out.P("  _ = ", out.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "Background()",
			GoImportPath: protogen.GoImportPath("context"),
		}))
		out.P("  _ = (*", out.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "Empty",
			GoImportPath: protogen.GoImportPath(s.ProtogenPackage),
		}), ")(nil)")
		out.P(")")
		out.P()

		implName := s.ImplName()
		out.P("type ", implName, " struct {")
		out.P("}")
		out.P()

		for _, rpc := range s.RPCs {
			if rpc.IsAsync() {
				out.P("func (s *", implName, ") ", rpc.Name, "(ctx context.Context, req *", s.ProtogenIdent(rpc.Req.GoName), ") {")
				out.P("}")
			} else {
				out.P("func (s *", implName, ") ", rpc.Name, "(ctx context.Context, req *", s.ProtogenIdent(rpc.Req.GoName), ", rsp *", s.ProtogenIdent(rpc.Rsp.GoName), ") error {")
				out.P("  return nil")
				out.P("}")
			}
			out.P()
		}
	}
}

// GenFromAst 通过ast生成impl
func (s *Service) GenFromAst(path string, out *protogen.GeneratedFile) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Panicf("parse file %s error: %v", path, err)
	}

	protogenPkg := filepath.Base(s.ProtogenPackage)

	// import至少需要包含 context 和 s.ProtogenPackage 两个包
	astutil.AddNamedImport(fset, f, "context", "context")
	astutil.AddNamedImport(fset, f, protogenPkg, s.ProtogenPackage)

	// 确认是否存在impl结构体
	implStruct := s.ImplName()
	if !hasStruct(f, implStruct) {
		// 添加结构体
		f.Decls = append(f.Decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(implStruct),
					Type: &ast.StructType{
						Fields: &ast.FieldList{},
					},
				},
			},
		})
	}

	// 检测是否包含所有的rpc方法
	allMethods := getMethodNamesForStruct(f, implStruct)
	var newFuncDecls []*ast.FuncDecl

	for _, rpc := range s.RPCs {
		if allMethods[rpc.Name] {
			continue // 这里就不进行签名校验了
		}
		log.Println("miss method:", implStruct, rpc.Name)
		var params = []*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("ctx")},
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("context"),
					Sel: ast.NewIdent("Context"),
				},
			},
			{
				Names: []*ast.Ident{ast.NewIdent("req")},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(protogenPkg),
						Sel: ast.NewIdent(rpc.Req.GoName),
					},
				},
			},
		}
		var results []*ast.Field
		var statements []ast.Stmt
		if rpc.IsAsync() {

		} else {
			params = append(params, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent("rsp")},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(protogenPkg),
						Sel: ast.NewIdent(rpc.Rsp.GoName),
					},
				},
			})
			results = []*ast.Field{
				{
					Type: ast.NewIdent("error"),
				},
			}
			statements = append(statements, &ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("nil"),
				},
			})
		}
		newFuncDecl := buildMethodFuncDecl(implStruct, rpc.Name, params, results, statements)
		newFuncDecls = append(newFuncDecls, newFuncDecl)
	}

	// 如果有新方法需要添加，使用特殊的处理方式
	if len(newFuncDecls) > 0 {
		genMixedOutput(fset, f, newFuncDecls, out)
		return
	}

	// 序列化并写入文件
	var buf = bytes.NewBuffer(nil)

	// 使用 printer.Fprint 而不是 format.Node，更好地保留注释
	err = printer.Fprint(buf, fset, f)
	if err != nil {
		log.Panicf("print file %s error: %v", path, err)
	}

	// 再进行格式化
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Panicf("format file %s error: %v", path, err)
	}
	out.Write(formatted)
}

func fileIsExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Panicf("file %s check error: %v", path, err)
		}
	}
	if info.IsDir() {
		log.Panicf("file %s is a directory", path)
	}
	return true
}

func hasStruct(f *ast.File, structName string) bool {
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if typeSpec.Name.Name != structName {
				continue
			}
			// 是 type Player struct {...}
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				return true
			}
		}
	}
	return false
}

func getMethodNamesForStruct(f *ast.File, structName string) map[string]bool {
	methods := map[string]bool{}
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}

		// 接收者表达式
		recvExpr := funcDecl.Recv.List[0].Type

		var recvIdent *ast.Ident

		switch expr := recvExpr.(type) {
		case *ast.Ident:
			recvIdent = expr
		case *ast.StarExpr: // 处理 *MyStruct
			if ident, ok := expr.X.(*ast.Ident); ok {
				recvIdent = ident
			}
		}

		if recvIdent != nil && recvIdent.Name == structName {
			methods[funcDecl.Name.Name] = true
		}
	}
	return methods
}

func buildMethodFuncDecl(
	structName string,
	methodName string,
	params []*ast.Field,
	results []*ast.Field,
	bodyStmts []ast.Stmt,
) *ast.FuncDecl {
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(methodName),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: ast.NewIdent(structName)},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: &ast.FieldList{List: results},
		},
		Body: &ast.BlockStmt{List: bodyStmts},
	}

	// 更可靠的注释设置方法
	comment := &ast.Comment{
		Slash: token.NoPos,
		Text:  "// " + methodName,
	}
	funcDecl.Doc = &ast.CommentGroup{
		List: []*ast.Comment{comment},
	}

	return funcDecl
}

// genMixedOutput 处理增量更新时的特殊情况，通过文本拼接避免AST位置信息不一致导致的注释丢失
func genMixedOutput(fset *token.FileSet, f *ast.File, newFuncDecls []*ast.FuncDecl, out *protogen.GeneratedFile) {
	// 第一步：生成现有文件的格式化代码
	var existingBuf = bytes.NewBuffer(nil)
	err := printer.Fprint(existingBuf, fset, f)
	if err != nil {
		log.Panicf("print existing file error: %v", err)
	}

	formatted, err := format.Source(existingBuf.Bytes())
	if err != nil {
		log.Panicf("format existing file error: %v", err)
	}
	existingCode := string(formatted)

	// 第二步：为每个新方法单独生成代码
	var newMethodsCode strings.Builder

	for _, funcDecl := range newFuncDecls {
		// 创建临时文件AST，只包含这一个函数
		tempFile := &ast.File{
			Name:  ast.NewIdent("main"), // 这个不重要
			Decls: []ast.Decl{funcDecl},
		}

		var tempBuf = bytes.NewBuffer(nil)
		tempFset := token.NewFileSet()

		err := printer.Fprint(tempBuf, tempFset, tempFile)
		if err != nil {
			log.Panicf("print new method error: %v", err)
		}

		// 提取函数声明部分（去掉package声明）
		tempCode := tempBuf.String()
		lines := strings.Split(tempCode, "\n")
		var funcLines []string
		inFunc := false

		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "//") ||
				strings.HasPrefix(strings.TrimSpace(line), "func") {
				inFunc = true
			}
			if inFunc && strings.TrimSpace(line) != "" {
				funcLines = append(funcLines, line)
			}
		}

		if len(funcLines) > 0 {
			newMethodsCode.WriteString("\n")
			newMethodsCode.WriteString(strings.Join(funcLines, "\n"))
			newMethodsCode.WriteString("\n")
		}
	}

	// 第三步：拼接现有代码和新方法
	finalCode := existingCode + newMethodsCode.String()

	// 第四步：最终格式化
	finalFormatted, err := format.Source([]byte(finalCode))
	if err != nil {
		log.Panicf("format final file error: %v", err)
	}

	out.Write(finalFormatted)
}
