package formatter

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document/order"
)

// FormatValue takes a decoded value and a format string, and returns the value formatted according
// to the specified format.
func FormatValue(v decoder.Value, format string) (any, error) {
	switch format {
	case "DEC":
		if v.Len > 64 {
			return nil, fmt.Errorf("DEC requires <= 64 bits, got %d", v.Len)
		}
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

// GetBit retrieves the value of a specific bit from a decoded value, considering the specified bit order.
// It returns an error if the bit index is out of range or if the bit order is unknown.
func GetBit(v decoder.Value, idx int, bitOrder order.BitOrder) (uint64, error) {
	if idx < 0 || int64(idx) >= v.Len {
		return 0, fmt.Errorf("bit index %d out of range for %q", idx, v.Name)
	}

	switch bitOrder {
	case order.MSB_FIRST:
		shift := v.Len - int64(idx) - 1
		return (v.UInt >> shift) & 1, nil

	case order.LSB_FIRST:
		return (v.UInt >> idx) & 1, nil

	default:
		return 0, fmt.Errorf("unknown bit order %q", bitOrder)
	}
}

// ConvertMappedValue converts a string representation of a value to its corresponding Go type.
// It supports boolean values ("true" and "false") and returns the original string for other values.
func ConvertMappedValue(s string) any {
	switch s {
	case "true":
		return true
	case "false":
		return false
	default:
		return s
	}
}
