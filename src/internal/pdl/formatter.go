package pdl

import (
	"encoding/hex"
	"fmt"
)

func FormatValue(v Value, format string) (any, error) {
	switch format {

	case "DEC":
		return v.UInt, nil

	case "HEX":
		return hex.EncodeToString(v.Bits), nil

	case "BIN":
		return fmt.Sprintf("%0*b", v.Len, v.UInt), nil

	case "BOOL":
		return v.UInt != 0, nil

	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}
