package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"github.com/creotiv/gofmtc/internal/rules"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gofmtc <file-to-modify.go>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	// Load the Go source file
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err)
		os.Exit(1)
	}

	// Traverse and modify the AST
	ast.Inspect(node, func(n ast.Node) bool {
		// for function calls
		if ce, ok := n.(*ast.CallExpr); ok {
			rules.VlogMsgCallExpr(ce)
			rules.FmtErrorfCallExpr(ce)
			rules.ErrorsNewCallExpr(ce)
		}

		return true
	})

	// Print the modified AST with gofmt config
	printer.Fprint(os.Stdout, fset, node)
}
