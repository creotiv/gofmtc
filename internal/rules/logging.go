package rules

import (
	"go/ast"
	"go/token"
	"strings"
)

func VlogMsgCallExpr(ce *ast.CallExpr) {
	// Check for selector expression
	selExpr, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	// Check if the method is 'Msg'
	if selExpr.Sel.Name != "Msg" {
		return
	}

	// Check if .Msg() is preceded by .Error(), .Info(), .Debug(), or .Warn()
	if !isPrecededByLoggingMethod(selExpr) {
		return
	}

	// Ensure there is at least one argument
	if len(ce.Args) == 0 {
		return
	}

	// Handle the first argument
	handleVlogMsgArgument(ce.Args[0])
}

func handleVlogMsgArgument(arg ast.Expr) {
	switch v := arg.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			modifyVlogMsgStringLiteral(v)
		}
	case *ast.BinaryExpr:
		// Recursively handle the left part of the binary expression
		handleVlogMsgArgument(v.X)
	}
}

func modifyVlogMsgStringLiteral(lit *ast.BasicLit) {
	// Modify the string value
	strValue := lit.Value
	if len(strValue) > 2 { // considering the quotes
		// Capitalize the first character inside the quotes
		modifiedStr := "\"" + strings.ToUpper(strValue[1:2]) + strValue[2:]
		lit.Value = modifiedStr
	}
}

func isPrecededByLoggingMethod(selExpr *ast.SelectorExpr) bool {
	current := selExpr.X
	for current != nil {
		switch x := current.(type) {
		case *ast.CallExpr:
			if se, ok := x.Fun.(*ast.SelectorExpr); ok {
				if isLoggingMethod(se.Sel.Name) {
					return true
				}
				current = se.X
			} else {
				return false
			}
		default:
			// If the current expression is not a call expression, end the search
			return false
		}
	}
	return false
}

func isLoggingMethod(methodName string) bool {
	switch methodName {
	case "Error", "Info", "Debug", "Warn":
		return true
	default:
		return false
	}
}

///////////////////////////////////////////////////////////////////////////////

func FmtErrorfCallExpr(ce *ast.CallExpr) {
	// Check if this is a call to fmt.Errorf
	if !isFmtErrorfCall(ce) {
		return
	}

	// Modify the first argument if it's a string literal or a binary expression
	if len(ce.Args) > 0 {
		handleFmtErrorfArgument(ce.Args[0])
	}
}

func isFmtErrorfCall(ce *ast.CallExpr) bool {
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return id.Name == "fmt" && sel.Sel.Name == "Errorf"
}

func handleFmtErrorfArgument(arg ast.Expr) {
	switch v := arg.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			modifyFmtErrorfStringLiteralStart(v)
			modifyFmtErrorfStringLiteralAppend(nil, v)
		}
	case *ast.BinaryExpr:
		handleFmtErrorfArgument(v.X)
		modifyFmtErrorfStringLiteralAppend(
			v.X,
			v.Y.(*ast.BasicLit))
	}
}

func modifyFmtErrorfStringLiteralStart(lit *ast.BasicLit) {
	// Lowercase the first letter
	strValue := lit.Value
	if len(strValue) > 2 {
		modifiedStr := "\"" + strings.ToLower(strValue[1:2]) + strValue[2:]
		lit.Value = modifiedStr
	}
}

func modifyFmtErrorfStringLiteralAppend(prev interface{}, lit *ast.BasicLit) {
	// add space before '%w'
	modifiedStr := lit.Value
	if len(modifiedStr) > 2 {
		modifiedStr := strings.Replace(modifiedStr, ":%w", ": %w", 1)
		lit.Value = modifiedStr
	}

	if prev != nil {
		var prevLit *ast.BasicLit

		switch v := prev.(type) {
		case *ast.BasicLit:
			prevLit = v
		case *ast.BinaryExpr:
			prevLit = v.Y.(*ast.BasicLit)
		}

		modifiedStr = prevLit.Value

		if len(modifiedStr) > 2 {
			modifiedStr := strings.Replace(modifiedStr, ":%w", ": %w", 1)
			prevLit.Value = modifiedStr
		}

		modifiedStr = lit.Value
		if string(prevLit.Value[len(prevLit.Value)-2]) == ":" &&
			string(lit.Value[1:3]) == "%w" {
			modifiedStr := strings.Replace(modifiedStr, "%w", " %w", 1)
			lit.Value = modifiedStr
		}
	}
}

//////////////////////////////////////////////////////////////////////////////
func ErrorsNewCallExpr(ce *ast.CallExpr) {
	// Check if this is a call to fmt.Errorf
	if !isErrorsNewCall(ce) {
		return
	}

	// Modify the first argument if it's a string literal or a binary expression
	if len(ce.Args) > 0 {
		handleErrorsNewArgument(ce.Args[0])
	}
}

func isErrorsNewCall(ce *ast.CallExpr) bool {
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return id.Name == "errors" && sel.Sel.Name == "New"
}

func handleErrorsNewArgument(arg ast.Expr) {
    switch v := arg.(type) {
    case *ast.BasicLit:
        if v.Kind == token.STRING {
            modifyErrorsNewStringLiteral(v)
        }
    case *ast.BinaryExpr:
        handleErrorsNewArgument(v.X)
    }
}

func modifyErrorsNewStringLiteral(lit *ast.BasicLit) {
    strValue := lit.Value
    if len(strValue) > 2 {
        // Lowercase the first letter
        modifiedStr := "\"" + strings.ToLower(strValue[1:2]) + strValue[2:]
        lit.Value = modifiedStr
    }
}
