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
		if arr, ok := result.Arrays[rule.Field]; ok {
			if rule.Format != "" && rule.Format != "ARRAY" {
				return nil, fmt.Errorf("array field %q only supports ARRAY format", rule.Field)
			}

			items := make([]any, 0, len(arr.Items))

			for i, item := range arr.Items {
				childDoc, ok := set.Documents[arr.Packet]
				if !ok {
					return nil, fmt.Errorf("unknown packet %q", arr.Packet)
				}

				childJSON, err := BuildJSONWithSet(set, childDoc, item.Result)
				if err != nil {
					return nil, fmt.Errorf("%w in %q[%d]", err, arr.Packet, i)
				}

				items = append(items, childJSON)
			}

			if err := setJSONPath(&res, rule.Path, items); err != nil {
				return nil, err
			}

			continue
		}

		value, ok := result.Values[rule.Field]
		if !ok {
			return nil, fmt.Errorf("output field %q is not decoded", rule.Field)
		}

		var outValue any

		if rule.UseAsSwitch {
			target, err := decoder.ResolveOutAsSwitch(root, result, rule)
			if err != nil {
				return nil, err
			}

			if formatter.IsSupportedFormat(target) {
				formatted, err := formatter.FormatValue(value, target)
				if err != nil {
					return nil, err
				}
				outValue = formatted
			} else {
				childDoc, ok := set.Documents[target]
				if !ok {
					return nil, fmt.Errorf("unknown packet %q", target)
				}

				childResult, err := decoder.DecodeWithSet(set, childDoc, value.Bits)
				if err != nil {
					return nil, fmt.Errorf("%w in %q", err, target)
				}

				childJSON, err := BuildJSONWithSet(set, childDoc, childResult)
				if err != nil {
					return nil, fmt.Errorf("%w in %q", err, target)
				}

				outValue = childJSON
			}
		} else if rule.AsPacket != "" {
			childDoc, ok := set.Documents[rule.AsPacket]
			if !ok {
				return nil, fmt.Errorf("unknown packet %q", rule.AsPacket)
			}

			childResult, err := decoder.DecodeWithSet(set, childDoc, value.Bits)
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
