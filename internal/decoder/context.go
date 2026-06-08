package decoder

import (
	"fmt"
	"strconv"

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
	if def.UseSwitch {
		if def.Switch == nil {
			return fmt.Errorf("decode %s: missing switch body", def.Name)
		}

		selector, err := c.evalExpr(def.Switch.Selector)
		if err != nil {
			return fmt.Errorf("decode %s switch selector: %w", def.Name, err)
		}

		key := strconv.FormatInt(selector, 10)

		layout, ok := def.Switch.Cases[key]
		if !ok {
			if def.Switch.Default == nil {
				return fmt.Errorf("decode %s: no switch case for %s", def.Name, key)
			}
			layout = *def.Switch.Default
		}

		return c.decodeLayout(def.Name, layout)
	}

	layout := document.DefLayout{
		From:      def.From,
		Length:    def.Length,
		To:        def.To,
		UseLength: def.UseLength,
		UseTo:     def.UseTo,
	}

	return c.decodeLayout(def.Name, layout)
}

// decodeLayout decodes a field based on the provided layout, which specifies how to determine the starting position and length of the field.
func (c *decodeContext) decodeLayout(name string, layout document.DefLayout) error {
	from, err := c.evalExpr(layout.From)
	if err != nil {
		return fmt.Errorf("decode %s from: %w", name, err)
	}

	var length int64

	if layout.UseLength {
		length, err = c.evalExpr(layout.Length)
		if err != nil {
			return fmt.Errorf("decode %s length: %w", name, err)
		}
	} else if layout.UseTo {
		to, err := c.evalExpr(layout.To)
		if err != nil {
			return fmt.Errorf("decode %s to: %w", name, err)
		}

		if _, ok := layout.To.(document.EndExpr); ok {
			length = int64(len(c.data))*8 - from
		} else {
			length = to - from + 1
		}
	} else {
		return fmt.Errorf("decode %s: missing length or to", name)
	}

	if from < 0 {
		return fmt.Errorf("decode %s: from is negative: %d", name, from)
	}
	if length < 0 {
		return fmt.Errorf("decode %s: length is negative: %d", name, length)
	}

	totalBits := int64(len(c.data)) * 8
	if from+length > totalBits {
		return fmt.Errorf(
			"decode %s: field exceeds packet size: from=%d length=%d total=%d",
			name, from, length, totalBits,
		)
	}

	bits, err := extractBits(c.data, from, length)
	if err != nil {
		return fmt.Errorf("decode %s: %w", name, err)
	}

	var u uint64
	if length <= 64 {
		u, err = bitsToUint(bits, length, c.doc.ByteOrder)
		if err != nil {
			return fmt.Errorf("decode %s: %w", name, err)
		}
	}

	c.values[name] = Value{
		Name: name,
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
