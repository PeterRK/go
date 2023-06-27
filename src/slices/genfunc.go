// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// This program is run via "go generate" (via a directive in sort_ordered.go)
// to generate sort_func.go.

package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
)

var hackedFuncs = make(map[string]bool)

func main() {
	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, "sort_ordered.go", nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	af.Doc = nil
	af.Imports = nil
	af.Comments = nil

	var newDecl []ast.Decl
	for _, d := range af.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || fd.Name.IsExported() ||
			fd.Type.TypeParams == nil || len(fd.Type.TypeParams.List) != 1 {
			continue
		}
		field := fd.Type.TypeParams.List[0]
		if expr, ok := field.Type.(*ast.SelectorExpr); !ok ||
			expr.Sel.Name != "Ordered" ||
			len(field.Names) != 1 || field.Names[0].Name != "E" {
			continue
		}
		hackedFuncs[fd.Name.Name] = true
		fd.Type.TypeParams = nil
		newDecl = append(newDecl, fd)
	}
	af.Decls = newDecl
	ast.Walk(visitFunc(rewriteCalls), af)

	var out bytes.Buffer
	if err := format.Node(&out, fset, af); err != nil {
		log.Fatalf("format.Node: %v", err)
	}
	tpl := out.Bytes()

	funcPtn := regexp.MustCompile(`\nfunc `)
	src := funcPtn.ReplaceAll(tpl, []byte("\nfunc (cmp compare[E]) "))
	dumpOrDie("sort_func.go", src)
}

type visitFunc func(ast.Node) ast.Visitor

func (f visitFunc) Visit(n ast.Node) ast.Visitor { return f(n) }

func rewriteCalls(n ast.Node) ast.Visitor {
	ce, ok := n.(*ast.CallExpr)
	if ok {
		ident, ok := ce.Fun.(*ast.Ident)
		if ok && hackedFuncs[ident.Name] {
			ident.Name = "cmp." + ident.Name
		}
	}
	return visitFunc(rewriteCalls)
}

var header = `// Code generated from sort.go using genzfunc.go; DO NOT EDIT.

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

`

func dumpOrDie(filename string, src []byte) {
	src, err := format.Source(src)
	if err != nil {
		log.Fatalf("format.Source: %v on\n%s", err, src)
	}
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	if _, err := out.WriteString(header); err != nil {
		log.Fatal(err)
	}
	if _, err := out.Write(src); err != nil {
		log.Fatal(err)
	}
}
