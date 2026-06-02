package decoder

import (
	"fmt"

	"github.com/elecbug/pdl/internal/document"
)

// decodeContext holds the state during the decoding process, including the document, input data,
// encoded values, and variables.
type decodeContext struct {
	// The document being decoded, containing definitions and variables.
	doc *document.Document
	// The input packet data to decode.
	data []byte
	// A map of decoded field values, keyed by field name.
	values map[string]Value
	// A map of variable values, keyed by variable name.
	vars map[string]int64
}

// decodeDef decodes a single field definition from the document, extracting the specified bits from
// the input data and storing the result in the context's values map.
func (c *decodeContext) decodeDef(def document.Def) error {
	from, err := c.evalExpr(def.From)
	if err != nil {
		return fmt.Errorf("decode %s from: %w", def.Name, err)
	}

	var length int64

	if def.UseLength {
		length, err = c.evalExpr(def.Length)
		if err != nil {
			return fmt.Errorf("decode %s length: %w", def.Name, err)
		}
	} else if def.UseTo {
		to, err := c.evalExpr(def.To)
		if err != nil {
			return fmt.Errorf("decode %s to: %w", def.Name, err)
		}

		if _, ok := def.To.(document.EndExpr); ok {
			length = int64(len(c.data))*8 - from
		} else {
			length = to - from + 1
		}
	} else {
		return fmt.Errorf("decode %s: missing length or to", def.Name)
	}

	if from < 0 {
		return fmt.Errorf("decode %s: from is negative: %d", def.Name, from)
	}
	if length < 0 {
		return fmt.Errorf("decode %s: length is negative: %d", def.Name, length)
	}

	totalBits := int64(len(c.data)) * 8
	if from+length > totalBits {
		return fmt.Errorf(
			"decode %s: field exceeds packet size: from=%d length=%d total=%d",
			def.Name, from, length, totalBits,
		)
	}

	bits, err := extractBits(c.data, from, length)
	if err != nil {
		return fmt.Errorf("decode %s: %w", def.Name, err)
	}

	var u uint64

	if length <= 64 {
		u, err = bitsToUint(bits, length, c.doc.ByteOrder)
		if err != nil {
			return fmt.Errorf("decode %s: %w", def.Name, err)
		}
	}

	c.values[def.Name] = Value{
		Name: def.Name,
		Bits: bits,
		Len:  length,
		UInt: u,
	}

	return nil
}

// evalExpr evaluates an expression in the context of the current variables and decoded values,
// returning the resulting integer value or an error if the expression is invalid.
func (c *decodeContext) evalExpr(expr document.Expr) (int64, error) {
	switch e := expr.(type) {
	case document.NumberExpr:
		return e.Value, nil

	case document.IdentExpr:
		v, ok := c.vars[e.Name]
		if !ok {
			return 0, fmt.Errorf("undefined variable %q", e.Name)
		}
		return v, nil

	case document.FieldValueExpr:
		v, ok := c.values[e.Name]
		if !ok {
			return 0, fmt.Errorf("field %q is not decoded yet", e.Name)
		}
		return int64(v.UInt), nil

	case document.EndExpr:
		return int64(len(c.data))*8 - 1, nil

	case document.BinaryExpr:
		left, err := c.evalExpr(e.Left)
		if err != nil {
			return 0, err
		}

		right, err := c.evalExpr(e.Right)
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
