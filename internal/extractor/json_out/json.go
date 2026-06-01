package json_out

import (
	"fmt"
	"strconv"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/formatter"
)

func BuildJSON(doc *document.Document, result *decoder.DecodeResult) (any, error) {
	root := map[string]any{}

	for _, rule := range doc.Outs {
		value, ok := result.Values[rule.Field]
		if !ok {
			return nil, fmt.Errorf("output field %q is not decoded", rule.Field)
		}

		var outValue any

		if rule.BitIndex != nil {
			bit, err := formatter.GetBit(value, *rule.BitIndex, doc.BitOrder)
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
