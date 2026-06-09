package decoder

import (
	"fmt"

	"github.com/elecbug/pdl/internal/document"
)

// EvalExpr evaluates an expression in the context of decoding, allowing access to variables and decoded field values.
func (c *decodeContext) evalExpr(expr document.Expr) (int64, error) {
	return c.evalExprWithExtra(expr, nil)
}

// evalBoolExpr evaluates a boolean expression in the context of decoding, allowing access to variables and decoded field values, and returns the result as a boolean.
func (c *decodeContext) evalBoolExpr(expr document.Expr, extra map[string]int64) (bool, error) {
	v, err := c.evalExprWithExtra(expr, extra)
	if err != nil {
		return false, err
	}
	return v != 0, nil
}

// evalExprWithExtra evaluates an expression in the context of decoding, allowing access to variables, decoded field values,
// and additional variables provided in the extra map. It supports number literals, identifiers, field value references,
// binary expressions, and returns an error for unknown expression types or operators.
func (c *decodeContext) evalExprWithExtra(expr document.Expr, extra map[string]int64) (int64, error) {
	return evalExprCore(
		expr,

		func(name string) (int64, bool) {
			if extra != nil {
				if v, ok := extra[name]; ok {
					return v, true
				}
			}

			v, ok := c.vars[name]
			return v, ok
		},

		func(name string) (int64, bool) {
			v, ok := c.values[name]
			if !ok {
				return 0, false
			}

			return int64(v.UInt), true
		},

		func() (int64, error) {
			return int64(len(c.data))*8 - 1, nil
		},
	)
}

// evalBinary evaluates a binary expression given the operator and the left and right operand values, supporting arithmetic, comparison, and logical operators.
func evalBinary(op string, left, right int64) (int64, error) {
	switch op {
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

	case "<":
		return boolToInt(left < right), nil
	case ">":
		return boolToInt(left > right), nil
	case "<=":
		return boolToInt(left <= right), nil
	case ">=":
		return boolToInt(left >= right), nil
	case "==":
		return boolToInt(left == right), nil
	case "!=":
		return boolToInt(left != right), nil
	case "&&":
		return boolToInt(left != 0 && right != 0), nil
	case "||":
		return boolToInt(left != 0 || right != 0), nil

	default:
		return 0, fmt.Errorf("unknown operator %q", op)
	}
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
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
		case "<":
			if left < right {
				return 1, nil
			}
			return 0, nil
		case ">":
			if left > right {
				return 1, nil
			}
			return 0, nil
		case "<=":
			if left <= right {
				return 1, nil
			}
			return 0, nil
		case ">=":
			if left >= right {
				return 1, nil
			}
			return 0, nil
		case "==":
			if left == right {
				return 1, nil
			}
			return 0, nil
		case "!=":
			if left != right {
				return 1, nil
			}
			return 0, nil
		case "&&":
			if left != 0 && right != 0 {
				return 1, nil
			}
			return 0, nil
		case "||":
			if left != 0 || right != 0 {
				return 1, nil
			}
			return 0, nil
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
func evalOutExpr(root *document.Document, result *Result, expr document.Expr, extra map[string]int64) (int64, error) {
	return evalExprCore(
		expr,

		func(name string) (int64, bool) {

			if extra != nil {
				if v, ok := extra[name]; ok {
					return v, true
				}
			}

			for _, v := range root.Vars {
				if v.Name == name {
					value, err := evalOutExpr(
						root,
						result,
						v.Expr,
						extra,
					)

					if err != nil {
						return 0, false
					}

					return value, true
				}
			}

			return 0, false
		},

		func(name string) (int64, bool) {
			v, ok := result.Values[name]
			if !ok {
				return 0, false
			}

			return int64(v.UInt), true
		},

		func() (int64, error) {
			return 0, fmt.Errorf("end is not valid in out expression")
		},
	)
}

// evalOutBoolExpr evaluates a boolean expression in the context of output formatting, allowing access to variables and decoded field values, and returns the result as a boolean.
func evalOutBoolExpr(root *document.Document, result *Result, expr document.Expr, extra map[string]int64) (bool, error) {
	v, err := evalOutExpr(root, result, expr, extra)
	if err != nil {
		return false, err
	}
	return v != 0, nil
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

// evalOutExpr evaluates an expression in the context of output formatting, allowing access to variables and decoded field values.
func evalExprCore(
	expr document.Expr,
	resolveVar func(string) (int64, bool),
	resolveField func(string) (int64, bool),
	resolveEnd func() (int64, error),
) (int64, error) {

	switch e := expr.(type) {

	case document.NumberExpr:
		return e.Value, nil

	case document.IdentExpr:
		v, ok := resolveVar(e.Name)
		if !ok {
			return 0, fmt.Errorf("undefined variable %q", e.Name)
		}
		return v, nil

	case document.FieldValueExpr:
		v, ok := resolveField(e.Name)
		if !ok {
			return 0, fmt.Errorf("field %q not found", e.Name)
		}
		return v, nil

	case document.EndExpr:
		return resolveEnd()

	case document.BinaryExpr:
		left, err := evalExprCore(
			e.Left,
			resolveVar,
			resolveField,
			resolveEnd,
		)
		if err != nil {
			return 0, err
		}

		right, err := evalExprCore(
			e.Right,
			resolveVar,
			resolveField,
			resolveEnd,
		)
		if err != nil {
			return 0, err
		}

		return evalBinary(e.Op, left, right)

	default:
		return 0, fmt.Errorf("unknown expression type %T", expr)
	}
}
