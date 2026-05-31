package extractor

import (
	"fmt"
	"strconv"

	"github.com/elecbug/pdl/internal/ast"
	"github.com/elecbug/pdl/internal/decoder"
)

func BuildJSON(doc *ast.Document, result *decoder.DecodeResult) (any, error) {
	root := map[string]any{}

	for _, rule := range doc.Outs {
		value, ok := result.Values[rule.Field]
		if !ok {
			return nil, fmt.Errorf("output field %q is not decoded", rule.Field)
		}

		var outValue any

		if rule.BitIndex != nil {
			bit, err := GetBit(value, *rule.BitIndex, doc.BitOrder)
			if err != nil {
				return nil, err
			}

			key := strconv.FormatUint(bit, 10)
			outValue = bit

			if rule.Map != nil {
				if mapped, ok := rule.Map[key]; ok {
					outValue = convertMappedValue(mapped)
				}
			}
		} else if rule.Map != nil {
			key := strconv.FormatUint(value.UInt, 10)
			outValue = value.UInt

			if mapped, ok := rule.Map[key]; ok {
				outValue = convertMappedValue(mapped)
			}
		} else {
			formatted, err := FormatValue(value, rule.Format)
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

func convertMappedValue(s string) any {
	switch s {
	case "true":
		return true
	case "false":
		return false
	default:
		return s
	}
}
