package pdl

import "fmt"

func BuildJSON(
	doc *Document,
	result *DecodeResult,
) (map[string]any, error) {

	out := make(map[string]any)

	for _, rule := range doc.Outputs {

		value := result.Values[rule.Field]

		if rule.BitIndex == nil {

			if len(rule.Map) == 0 {

				formatted, err := FormatValue(
					value,
					rule.Format,
				)
				if err != nil {
					return nil, err
				}

				out[rule.Path] = formatted
				continue
			}
		}

		if rule.BitIndex != nil {

			bit, err := GetBit(
				value,
				*rule.BitIndex,
			)
			if err != nil {
				return nil, err
			}

			key := fmt.Sprintf("%d", bit)

			if mapped, ok := rule.Map[key]; ok {

				switch mapped {
				case "true":
					out[rule.Path] = true
				case "false":
					out[rule.Path] = false
				default:
					out[rule.Path] = mapped
				}

			} else {
				out[rule.Path] = bit
			}
		}
	}

	return out, nil
}
