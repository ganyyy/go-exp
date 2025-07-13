package main

import (
	"go/token"
	"log"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/guess"
	"github.com/dave/dst/dstutil"
)

func TestDave3(t *testing.T) {
	fset := token.NewFileSet()
	// 初始化 decorator，开启 import 管理
	dec := decorator.NewDecoratorWithImports(fset, "main", goast.New())

	// —————— 步骤1：解析源码 A，提取 snippet stms ——————
	srcA := `
package tmp
import (
	"log"
	"time"

	"mypkg/abc"
)


func _() {
    t := time.Now()
    log.Printf("time now: %v, %d", t, abc.SomeFunc())
}
`
	tmpF, err := dec.Parse([]byte(srcA))
	if err != nil {
		log.Fatal(err)
	}
	snippet := tmpF.Decls[1].(*dst.FuncDecl).Body.List
	// snippet 中的 dst.Ident 已包含 Path: "time" 和 "log"

	// —————— 步骤2：解析源码 B，定位目标函数 ——————
	srcB := `
package main

func Hello() {
    println("hello")
}
`
	fB, err := dec.Parse([]byte(srcB))
	if err != nil {
		log.Fatal(err)
	}

	// —————— 步骤3：插入 snippet，每次深度 clone ——————
	dstutil.Apply(fB, func(c *dstutil.Cursor) bool {
		if fn, ok := c.Node().(*dst.FuncDecl); ok && fn.Name.Name == "Hello" {
			newList := []dst.Stmt{}
			for _, stmt := range snippet {
				stmt = dst.Clone(stmt).(dst.Stmt)
				newList = append(newList, stmt)
			}
			fn.Body.List = append(newList, fn.Body.List...)
		}
		return true
	}, nil)

	// —————— 步骤4：还原 AST，自动补 import ——————
	res := decorator.NewRestorerWithImports("main", guess.New())
	if err := res.Print(fB); err != nil {
		log.Fatal(err)
	}
}
