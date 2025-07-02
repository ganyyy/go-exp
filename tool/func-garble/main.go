package main

import (
	"flag"
	"fmt"
	"go/ast"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "input", "", "Path to the input file")
	flag.Parse()

	const insertBody = `
	
	nowt := time.Now()
	checkDeadline := nowt.Add(-time.Minute).UnixNano()
	current := lastCheck.Load()
	if current >= checkDeadline {
		goto next
	}
	if !lastCheck.CompareAndSwap(current, nowt.UnixNano()) {
		goto next
	}
	if a != 100 {
		fmt.Println("a is not 100")
	}
	if b != 200 {
		fmt.Println("b is not 200")
	}
next:

	`

	if inputFile == "" {
		panic("Input file must be specified")
	}

	insertStmts := generateStmts(insertBody)

	source, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	f, err := decorator.Parse(source)
	if err != nil {
		panic(err)
	}

	dstutil.Apply(f, func(c *dstutil.Cursor) bool {
		fn, ok := c.Node().(*dst.FuncDecl)
		if !ok || fn == nil || fn.Name == nil || fn.Body == nil {
			return true
		}

		name := fn.Name.Name
		if !ast.IsExported(name) {
			// Skip non-exported functions
			return true
		}

		var tmpInserts = make([]dst.Stmt, 0, len(insertStmts))
		for _, stmt := range insertStmts {
			// Create a copy of the statement to avoid modifying the original
			tmpInserts = append(tmpInserts, dst.Clone(stmt).(dst.Stmt))
		}

		fn.Body.List = append(tmpInserts, fn.Body.List...)
		return true
	}, nil)

	decorator.Fprint(os.Stdout, f)
}

func generateStmts(src string) []dst.Stmt {
	src = fmt.Sprintf(
		`
package tmp

func _() {		
		%s
}
		`, src,
	)

	tmpF, err := decorator.Parse(src)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse source: %v", err))
	}
	if len(tmpF.Decls) == 0 {
		panic("No declarations found in the temporary file")
	}
	fnDecl, ok := tmpF.Decls[0].(*dst.FuncDecl)
	if !ok {
		panic("First declaration is not a function declaration")
	}
	if fnDecl.Body == nil {
		panic("Function declaration has no body")
	}
	return fnDecl.Body.List
}
