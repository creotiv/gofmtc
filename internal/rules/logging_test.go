package rules_test

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
	"testing"

	"github.com/creotiv/gofmtc/internal/rules"

	"github.com/stretchr/testify/require"
)

func wrapCode(code string) string {
	return fmt.Sprintf(`package main

func main() {
	%s
}%s`, code, "\n")
}

func callRule(code string, f func(*ast.CallExpr)) string {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(
		fset, "", strings.NewReader(code), parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Traverse and modify the AST
	ast.Inspect(node, func(n ast.Node) bool {
		// for function calls
		if ce, ok := n.(*ast.CallExpr); ok {
			f(ce)
		}

		return true
	})

	var buf bytes.Buffer

	printer.Fprint(&buf, fset, node)

	return buf.String()
}

func TestVlogFormatting(t *testing.T) {
	r := require.New(t)

	code := wrapCode(`vlog.Info().Stack().Msg("hello world")`)
	expected := wrapCode(`vlog.Info().Stack().Msg("Hello world")`)

	formated := callRule(code, rules.VlogMsgCallExpr)
	r.Equal(expected, formated)

	code = wrapCode(`vlog.Info().
		Stack().
		Msg(
			"h" +
				"ello world" + "a" + "b",
		)`)
	expected = wrapCode(`vlog.Info().
		Stack().
		Msg(
			"H" +
				"ello world" + "a" + "b",
		)`)

	formated = callRule(code, rules.VlogMsgCallExpr)
	r.Equal(expected, formated)

	code = wrapCode(`Msg("hello world")`)
	expected = wrapCode(`Msg("hello world")`)

	formated = callRule(code, rules.VlogMsgCallExpr)
	r.Equal(expected, formated)
}

func TestFmtErrorfFormatting(t *testing.T) {
	r := require.New(t)

	code := wrapCode(
		`fmt.Println(fmt.Errorf("Hello world:%w", errors.New("FF")))`)
	expected := wrapCode(
		`fmt.Println(fmt.Errorf("hello world: %w", errors.New("FF")))`)

	formated := callRule(code, rules.FmtErrorfCallExpr)
	r.Equal(expected, formated)

	code = wrapCode(`fmt.Println(
		fmt.Errorf(
			"H"+
				"ello world:"+"%w",
			errors.New(
				"FF"),
		),
	)`)
	expected = wrapCode(`fmt.Println(
		fmt.Errorf(
			"h"+
				"ello world:"+" %w",
			errors.New(
				"FF"),
		),
	)`)

	formated = callRule(code, rules.FmtErrorfCallExpr)
	r.Equal(expected, formated)
}

func TestErrorsNewFormatting(t *testing.T) {
	r := require.New(t)

	code := wrapCode(
		`fmt.Println(fmt.Errorf("Hello world:%w", errors.New("fF")))`)
	expected := wrapCode(
		`fmt.Println(fmt.Errorf("Hello world:%w", errors.New("fF")))`)

	formated := callRule(code, rules.ErrorsNewCallExpr)
	r.Equal(expected, formated)

	code = wrapCode(`fmt.Println(
		fmt.Errorf(
			"H"+
				"ello world:"+"%w",
			errors.New(
				"fF"),
		),
	)`)
	expected = wrapCode(`fmt.Println(
		fmt.Errorf(
			"H"+
				"ello world:"+"%w",
			errors.New(
				"fF"),
		),
	)`)

	formated = callRule(code, rules.ErrorsNewCallExpr)
	r.Equal(expected, formated)
}
