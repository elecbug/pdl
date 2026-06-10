package decoder

import (
	"fmt"

	"github.com/elecbug/pdl/internal/document"
)

// decodeContext holds the state during the decoding process, including the document, input data,
// encoded values, and variables.
type decodeContext struct {
	// The set of documents available for decoding, which may include multiple packet definitions and a designated root document.
	set *document.DocumentSet
	// The document being decoded, containing definitions and variables.
	doc *document.Document
	// The input packet data to decode.
	data []byte
	// A map of decoded field values, keyed by field name.
	values map[string]Value
	// A map of variable values, keyed by variable name.
	vars map[string]int64
	// A map of decoded array values, keyed by field name.
	arrays map[string]ArrayValue
	// The total number of bits consumed during decoding, used for tracking the current position in the input data.
	consumedBits int64
}

// decodeDef decodes a single field definition from the document, extracting the specified bits from
// the input data and storing the result in the context's values map.
func (c *decodeContext) decodeDef(def document.Def) error {
	if def.UseArray {
		return c.decodeArray(def.Name, def.Array)
	}

	if def.UseSwitch {
		layout, ok, err := c.resolveDefSwitchLayout(def)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("decode %s: no matching switch case", def.Name)
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

	end := from + length
	if end > c.consumedBits {
		c.consumedBits = end
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

	var extra map[string]int64

	if rule.AsSwitch.Selector != nil {
		selector, err := evalOutExpr(root, result, rule.AsSwitch.Selector, nil)
		if err != nil {
			return "", fmt.Errorf("as switch selector for field %q: %w", rule.Field, err)
		}

		extra = map[string]int64{
			"val": selector,
		}
	}

	for _, cs := range rule.AsSwitch.Cases {
		ok, err := evalOutBoolExpr(root, result, cs.Condition, extra)
		if err != nil {
			return "", fmt.Errorf("as switch case for field %q: %w", rule.Field, err)
		}

		if ok {
			return cs.Value, nil
		}
	}

	if rule.AsSwitch.Default != nil {
		return *rule.AsSwitch.Default, nil
	}

	return "", fmt.Errorf("no matching as switch case for field %q", rule.Field)
}

// resolveDefSwitchLayout resolves the layout for a field definition that uses a switch statement by evaluating the selector
// expression and finding the matching case in the definition's Switch field. It returns the resolved layout, a boolean
// indicating whether a matching case was found, and an error if any issues occur during evaluation.
func (c *decodeContext) resolveDefSwitchLayout(def document.Def) (document.DefLayout, bool, error) {
	if def.Switch == nil {
		return document.DefLayout{}, false, fmt.Errorf("decode %s: missing switch body", def.Name)
	}

	var extra map[string]int64

	if def.Switch.Selector != nil {
		selector, err := c.evalExpr(def.Switch.Selector)
		if err != nil {
			return document.DefLayout{}, false, fmt.Errorf("decode %s switch selector: %w", def.Name, err)
		}

		extra = map[string]int64{
			"val": selector,
		}
	}

	for _, cs := range def.Switch.Cases {
		ok, err := c.evalBoolExpr(cs.Condition, extra)
		if err != nil {
			return document.DefLayout{}, false, fmt.Errorf("decode %s switch case: %w", def.Name, err)
		}

		if ok {
			return cs.Layout, true, nil
		}
	}

	if def.Switch.Default != nil {
		return *def.Switch.Default, true, nil
	}

	return document.DefLayout{}, false, nil
}

// decodeArray decodes an array of fields based on the provided DefArray structure, which specifies
// how to determine the starting position and count of the fields in the array.
func (c *decodeContext) decodeArray(name string, arr *document.DefArray) error {
	if arr == nil {
		return fmt.Errorf("decode %s: missing array body", name)
	}

	if c.set == nil {
		return fmt.Errorf("decode %s: array requires DocumentSet", name)
	}

	childDoc, ok := c.set.Documents[arr.Packet]
	if !ok {
		return fmt.Errorf("decode %s: unknown packet %q", name, arr.Packet)
	}

	from, err := c.evalExpr(arr.From)
	if err != nil {
		return fmt.Errorf("decode %s from: %w", name, err)
	}

	if from < 0 {
		return fmt.Errorf("decode %s: from is negative: %d", name, from)
	}

	totalBits := int64(len(c.data)) * 8
	if from > totalBits {
		return fmt.Errorf("decode %s: from exceeds packet size: from=%d total=%d", name, from, totalBits)
	}

	var maxCount int64
	if arr.CountToEnd {
		maxCount = -1
	} else {
		maxCount, err = c.evalExpr(arr.Count)
		if err != nil {
			return fmt.Errorf("decode %s count: %w", name, err)
		}
		if maxCount < 0 {
			return fmt.Errorf("decode %s: count is negative: %d", name, maxCount)
		}
	}

	out := ArrayValue{
		Name:   name,
		Packet: arr.Packet,
	}

	offset := from
	for i := int64(0); ; i++ {
		if arr.CountToEnd {
			if offset >= totalBits {
				break
			}
		} else {
			if i >= maxCount {
				break
			}
			if offset >= totalBits {
				return fmt.Errorf("decode %s[%d]: packet ended before count was satisfied", name, i)
			}
		}

		remain := totalBits - offset
		childBits, err := extractBits(c.data, offset, remain)
		if err != nil {
			return fmt.Errorf("decode %s[%d]: %w", name, i, err)
		}

		childResult, err := DecodeWithSet(c.set, childDoc, childBits)
		if err != nil {
			return fmt.Errorf("decode %s[%d] as %q: %w", name, i, arr.Packet, err)
		}

		if childResult.ConsumedBits <= 0 {
			return fmt.Errorf("decode %s[%d] as %q: consumed zero bits", name, i, arr.Packet)
		}

		itemBits, err := extractBits(c.data, offset, childResult.ConsumedBits)
		if err != nil {
			return fmt.Errorf("decode %s[%d] item bits: %w", name, i, err)
		}

		out.Items = append(out.Items, ArrayItem{
			Bits:   itemBits,
			Len:    childResult.ConsumedBits,
			Result: childResult,
		})

		offset += childResult.ConsumedBits
	}

	c.arrays[name] = out

	if offset > c.consumedBits {
		c.consumedBits = offset
	}

	return nil
}
