package formatter

import (
	"encoding/hex"
	"fmt"
	"net"
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
		return formatHex(v)

	case "BIN":
		return "0b" + fmt.Sprintf("%0*b", v.Len, v.UInt), nil

	case "BOOL":
		if v.Len != 1 {
			return nil, fmt.Errorf("BOOL requires 1 bit, got %d bits", v.Len)
		}
		return v.UInt != 0, nil

	case "ASCII":
		return formatASCII(v)
	case "UTF8":
		return formatUTF8(v)
	case "IP4":
		return formatIP4(v)
	case "IP6":
		return formatIP6(v)
	case "MAC":
		return formatMAC(v)
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

func formatHex(v decoder.Value) (string, error) {
	if v.Len == 0 {
		return "", nil
	}

	hexDigits := int((v.Len + 3) / 4)

	if v.Len <= 64 {
		return "0x" + strings.ToUpper(fmt.Sprintf("%0*X", hexDigits, v.UInt)), nil
	}

	s := strings.ToUpper(hex.EncodeToString(v.Bits))

	if len(s) > hexDigits {
		s = s[len(s)-hexDigits:]
	}

	if len(s) < hexDigits {
		s = strings.Repeat("0", hexDigits-len(s)) + s
	}

	return "0x" + s, nil
}

// formatASCII converts a decoded value to its ASCII string representation. It checks that the value is byte-aligned
// and returns an error if it is not. The function uses the BitsToBytes helper to convert the bit slice to a byte slice
// before converting it to a string.
func formatASCII(v decoder.Value) (string, error) {
	if v.Len%8 != 0 {
		return "", fmt.Errorf("ASCII requires byte-aligned data")
	}

	return string(v.Bits), nil
}

// formatUTF8 converts a decoded value to its UTF-8 string representation. Similar to formatASCII, it checks that the
// value is byte-aligned
func formatUTF8(v decoder.Value) (string, error) {
	if v.Len%8 != 0 {
		return "", fmt.Errorf("UTF8 requires byte-aligned data")
	}

	return string(v.Bits), nil
}

// formatIP4 converts a decoded value to its IPv4 string representation. It checks that the value is exactly 32 bits long
// and returns an error if it is not. The function uses the BitsToBytes helper to convert the bit slice to a byte slice
// before constructing the net.IP object and returning its string representation.
func formatIP4(v decoder.Value) (string, error) {
	if v.Len != 32 {
		return "", fmt.Errorf("IP4 requires 32 bits")
	}

	return net.IPv4(
		v.Bits[0],
		v.Bits[1],
		v.Bits[2],
		v.Bits[3],
	).String(), nil
}

// formatIP6 converts a decoded value to its IPv6 string representation. It checks that the value is exactly 128 bits long
// and returns an error if it is not. The function uses the BitsToBytes helper to convert the bit slice to a byte slice
// before constructing the net.IP object and returning its string representation.
func formatIP6(v decoder.Value) (string, error) {
	if v.Len != 128 {
		return "", fmt.Errorf("IP6 requires 128 bits")
	}

	ip := net.IP(v.Bits)

	return ip.String(), nil
}

// formatMAC converts a decoded value to its MAC address string representation. It checks that the value is exactly 48 bits long
// and returns an error if it is not. The function uses the BitsToBytes helper to convert the bit slice to a byte slice
// before formatting it as a MAC address string.
func formatMAC(v decoder.Value) (string, error) {
	if v.Len != 48 {
		return "", fmt.Errorf("MAC requires 48 bits")
	}

	return fmt.Sprintf(
		"%02X:%02X:%02X:%02X:%02X:%02X",
		v.Bits[0], v.Bits[1], v.Bits[2],
		v.Bits[3], v.Bits[4], v.Bits[5],
	), nil
}
