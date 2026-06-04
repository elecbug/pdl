package json_out

import (
	"fmt"
	"strconv"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/formatter"
)

// BuildJSON constructs a JSON-compatible Go value from decoded values using
// the document's output rules.
//
// It iterates over each output rule, retrieves the decoded field value,
// applies optional bit extraction, mapping, or formatting, and then writes
// the value at the configured JSON path.
func BuildJSON(doc *document.Document, result *decoder.Result) (any, error) {
	root := map[string]any{}

	for _, rule := range doc.Outs {
		value, ok := result.Values[rule.Field]
		if !ok {
			return nil, fmt.Errorf("output field %q is not decoded", rule.Field)
		}

		var outValue any

		if rule.HasBitIndex {
			bit, err := formatter.GetBit(value, rule.BitIndex, doc.BitOrder)
			if err != nil {
				return nil, err
			}

			key := strconv.FormatUint(bit, 10)
			outValue = bit

			if rule.Map != nil {
				if mapped, ok := rule.Map[key]; ok {
					outValue = formatter.ConvertMappedValue(mapped)
				}
			}
		} else if rule.Map != nil {
			key := strconv.FormatUint(value.UInt, 10)
			outValue = value.UInt

			if mapped, ok := rule.Map[key]; ok {
				outValue = formatter.ConvertMappedValue(mapped)
			}
		} else {
			formatted, err := formatter.FormatValue(value, rule.Format)
			if err != nil {
				return nil, err
			}
			outValue = formatted
		}

		if err := setJSONPath(&root, rule.Path, outValue); err != nil {
			return nil, err
		}
	}

	return root, nil
}

// BuildJSONWithSet is an extended version of BuildJSON that takes a DocumentSet as input, allowing it
// to handle nested packet definitions when building the JSON output. It checks for the AsPacket field
// in the output rules and recursively decodes and builds JSON for nested packets as needed.
func BuildJSONWithSet(set *document.DocumentSet, doc *document.Document, result *decoder.Result) (any, error) {
	root := map[string]any{}

	for _, rule := range doc.Outs {
		value, ok := result.Values[rule.Field]
		if !ok {
			return nil, fmt.Errorf("output field %q is not decoded", rule.Field)
		}

		var outValue any

		if rule.AsPacket != "" {
			childDoc, ok := set.Documents[rule.AsPacket]
			if !ok {
				return nil, fmt.Errorf("unknown packet %q", rule.AsPacket)
			}

			childResult, err := decoder.Decode(childDoc, value.Bits)
			if err != nil {
				return nil, fmt.Errorf("%w in %q", err, rule.AsPacket)
			}

			childJSON, err := BuildJSONWithSet(set, childDoc, childResult)
			if err != nil {
				return nil, fmt.Errorf("%w in %q", err, rule.AsPacket)
			}

			outValue = childJSON
		} else if rule.HasBitIndex {
			bit, err := formatter.GetBit(value, rule.BitIndex, doc.BitOrder)
			if err != nil {
				return nil, err
			}

			key := strconv.FormatUint(bit, 10)
			outValue = bit

			if rule.Map != nil {
				if mapped, ok := rule.Map[key]; ok {
					outValue = formatter.ConvertMappedValue(mapped)
				}
			}
		} else if rule.Map != nil {
			key := strconv.FormatUint(value.UInt, 10)
			outValue = value.UInt

			if mapped, ok := rule.Map[key]; ok {
				outValue = formatter.ConvertMappedValue(mapped)
			}
		} else {
			formatted, err := formatter.FormatValue(value, rule.Format)
			if err != nil {
				return nil, err
			}
			outValue = formatted
		}

		if err := setJSONPath(&root, rule.Path, outValue); err != nil {
			return nil, err
		}
	}

	return root, nil
}
