package pdl

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func FormatValue(v Value, format string) (any, error) {
	switch format {
	case "DEC":
		return v.UInt, nil

	case "HEX":
		return strings.ToUpper(hex.EncodeToString(v.Bits)), nil

	case "BIN":
		return fmt.Sprintf("%0*b", v.Len, v.UInt), nil

	case "BOOL":
		if v.Len != 1 {
			return nil, fmt.Errorf("BOOL requires 1 bit, got %d bits", v.Len)
		}
		return v.UInt != 0, nil

	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}

func GetBit(v Value, idx int, mode string) (uint64, error) {
	if idx < 0 || int64(idx) >= v.Len {
		return 0, fmt.Errorf("bit index %d out of range for %q", idx, v.Name)
	}

	switch mode {
	case "", "BIG_ENDIAN":
		shift := v.Len - int64(idx) - 1
		return (v.UInt >> shift) & 1, nil

	case "LITTLE_ENDIAN":
		return (v.UInt >> idx) & 1, nil

	default:
		return 0, fmt.Errorf("unknown mode %q", mode)
	}
}
