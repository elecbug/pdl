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

// ResolveOutAsSwitch resolves the output field name based on the value of a selector expression in an "as switch" rule.
// It evaluates the selector expression and looks up the corresponding case in the rule's AsSwitch mapping, returning the
// target field name or an error if no matching case is found.
func ResolveOutAsSwitch(root *document.Document, result *Result, rule document.Out) (string, error) {
	if rule.AsSwitch == nil {
		return "", fmt.Errorf("missing as switch body for field %q", rule.Field)
	}

	selector, err := evalOutExpr(root, result, rule.AsSwitch.Selector)
	if err != nil {
		return "", fmt.Errorf("as switch selector for field %q: %w", rule.Field, err)
	}

	key := strconv.FormatInt(selector, 10)

	if target, ok := rule.AsSwitch.Cases[key]; ok {
		return target, nil
	}

	if rule.AsSwitch.Default != nil {
		return *rule.AsSwitch.Default, nil
	}

	return "", fmt.Errorf("no as switch case for field %q selector=%s", rule.Field, key)
}
