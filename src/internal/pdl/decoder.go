package pdl

import (
	"fmt"
)

type Value struct {
	Name string

	Bits []byte
	Len  int64

	UInt uint64

	Mode string
}

type DecodeResult struct {
	Values map[string]Value
}

func Decode(doc *Document, data []byte) (*DecodeResult, error) {
	ctx := &decodeContext{
		doc:    doc,
		data:   data,
		values: make(map[string]Value),
	}

	for _, def := range doc.Defs {
		if err := ctx.decodeDef(def); err != nil {
			return nil, err
		}
	}

	return &DecodeResult{
		Values: ctx.values,
	}, nil
}

type decodeContext struct {
	doc    *Document
	data   []byte
	values map[string]Value
}

func (c *decodeContext) decodeDef(def Def) error {
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

		if _, ok := def.To.(EndExpr); ok {
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

	bits := extractBits(c.data, from, length)

	u, err := bitsToUint(bits, length, c.doc.Mode)
	if err != nil {
		return fmt.Errorf("decode %s: %w", def.Name, err)
	}

	c.values[def.Name] = Value{
		Name: def.Name,
		Bits: bits,
		Len:  length,
		UInt: u,
	}

	return nil
}

func (c *decodeContext) evalExpr(expr Expr) (int64, error) {
	switch e := expr.(type) {
	case NumberExpr:
		return e.Value, nil

	case IdentExpr:
		return 0, fmt.Errorf("identifier expression is not supported yet: %s", e.Name)

	case FieldValueExpr:
		v, ok := c.values[e.Name]
		if !ok {
			return 0, fmt.Errorf("field %q is not decoded yet", e.Name)
		}
		return int64(v.UInt), nil

	case EndExpr:
		return int64(len(c.data))*8 - 1, nil

	case BinaryExpr:
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
