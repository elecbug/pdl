package decoder

import (
	"fmt"

	"github.com/elecbug/pdl/internal/document"
)

// EvalExpr evaluates an expression in the context of decoding, allowing access to variables and decoded field values.
func (c *decodeContext) evalExpr(expr document.Expr) (int64, error) {
	return evalExprWithEnv(c, expr)
}

// evalExprWithEnv evaluates an expression in the context of the provided environment, allowing access to variables, decoded
// field values, and the "end" expression. It supports number literals, identifiers, field value references, binary expressions,
// and returns an error for unknown expression types or operators.
func evalExprWithEnv(env exprEnv, expr document.Expr) (int64, error) {
	switch e := expr.(type) {
	case document.NumberExpr:
		return e.Value, nil

	case document.IdentExpr:
		return env.ResolveIdent(e.Name)

	case document.FieldValueExpr:
		return env.ResolveField(e.Name)

	case document.EndExpr:
		return env.ResolveEnd()

	case document.BinaryExpr:
		left, err := evalExprWithEnv(env, e.Left)
		if err != nil {
			return 0, err
		}

		right, err := evalExprWithEnv(env, e.Right)
		if err != nil {
			return 0, err
		}

		switch e.Op {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("unknown operator %q", e.Op)
		}

	default:
		return 0, fmt.Errorf("unknown expression type %T", expr)
	}
}

// evalOutExpr evaluates an expression in the context of output formatting, allowing access to variables and decoded field values.
func evalOutExpr(root *document.Document, result *Result, expr document.Expr) (int64, error) {
	return evalExprWithEnv(&outExprEnv{
		root:   root,
		result: result,
	}, expr)
}

// exprEnv defines an interface for resolving identifiers, field values, and the "end" expression during expression evaluation.
type exprEnv interface {
	// ResolveIdent resolves a variable name to its value in the context of decoding, returning an error if the variable is not defined.
	ResolveIdent(name string) (int64, error)
	// ResolveField resolves a field name to its decoded unsigned integer value in the context of decoding, returning an error if the field is not decoded yet.
	ResolveField(name string) (int64, error)
	// ResolveEnd resolves the "end" expression to the total number of bits in the input data minus one, allowing it to be used as a dynamic endpoint in field definitions.
	ResolveEnd() (int64, error)
}

// ResolveIdent resolves a variable name to its value in the context of decoding, returning an error if the variable is not defined.
func (c *decodeContext) ResolveIdent(name string) (int64, error) {
	v, ok := c.vars[name]
	if !ok {
		return 0, fmt.Errorf("undefined variable %q", name)
	}
	return v, nil
}

// ResolveField resolves a field name to its decoded unsigned integer value in the context of decoding, returning an error if the field is not decoded yet.
func (c *decodeContext) ResolveField(name string) (int64, error) {
	v, ok := c.values[name]
	if !ok {
		return 0, fmt.Errorf("field %q is not decoded yet", name)
	}
	return int64(v.UInt), nil
}

// ResolveEnd resolves the "end" expression to the total number of bits in the input data minus one, allowing it to be used as a dynamic endpoint in field definitions.
func (c *decodeContext) ResolveEnd() (int64, error) {
	return int64(len(c.data))*8 - 1, nil
}

// outExprEnv implements the exprEnv interface for evaluating expressions in the context of output formatting, allowing access to
// variables and decoded field values from the document and result.
type outExprEnv struct {
	// The document containing variable definitions and byte order information.
	root *document.Document
	// The result containing decoded field values.
	result *Result
}

// ResolveIdent resolves a variable name to its value in the context of output formatting, returning an error if the variable is not defined.
func (e *outExprEnv) ResolveIdent(name string) (int64, error) {
	for _, v := range e.root.Vars {
		if v.Name == name {
			return evalExprWithEnv(e, v.Expr)
		}
	}
	return 0, fmt.Errorf("undefined variable %q", name)
}

// ResolveField resolves a field name to its decoded unsigned integer value in the context of output formatting, returning an error if the field is not decoded yet.
func (e *outExprEnv) ResolveField(name string) (int64, error) {
	v, ok := e.result.Values[name]
	if !ok {
		return 0, fmt.Errorf("field %q is not decoded yet", name)
	}
	return int64(v.UInt), nil
}

// ResolveEnd resolves the "end" expression in the context of output formatting, returning an error since "end" is not valid in this context.
func (e *outExprEnv) ResolveEnd() (int64, error) {
	return 0, fmt.Errorf("end is not valid in out as switch selector")
}
