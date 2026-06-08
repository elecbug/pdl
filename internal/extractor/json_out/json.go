package json_out

import (
	"fmt"
	"strconv"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/formatter"
)

// BuildJSONWithSet is an extended version of BuildJSON that takes a DocumentSet as input, allowing it
// to handle nested packet definitions when building the JSON output. It checks for the AsPacket field
// in the output rules and recursively decodes and builds JSON for nested packets as needed.
func BuildJSONWithSet(set *document.DocumentSet, root *document.Document, result *decoder.Result) (any, error) {
	res := map[string]any{}

	for _, rule := range root.Outs {
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
			bit, err := formatter.GetBit(value, rule.BitIndex, root.BitOrder)
			if err != nil {
				return nil, err
			}

			key := strconv.FormatUint(bit, 10)
			outValue = bit

			if rule.Map != nil {
				if mapped, ok := rule.Map[key]; ok {
					outValue = formatter.ConvertMappedValue(mapped)
				} else if rule.MapDefault != nil {
					outValue = formatter.ConvertMappedValue(*rule.MapDefault)
				}
			}
		} else if rule.Map != nil {
			key := strconv.FormatUint(value.UInt, 10)
			outValue = value.UInt

			if mapped, ok := rule.Map[key]; ok {
				outValue = formatter.ConvertMappedValue(mapped)
			} else if rule.MapDefault != nil {
				outValue = formatter.ConvertMappedValue(*rule.MapDefault)
			}
		} else {
			formatted, err := formatter.FormatValue(value, rule.Format)
			if err != nil {
				return nil, err
			}
			outValue = formatted
		}

		if err := setJSONPath(&res, rule.Path, outValue); err != nil {
			return nil, err
		}
	}

	return res, nil
}
