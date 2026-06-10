package pdl

import "github.com/elecbug/pdl/internal/token"

// PayloadFormat is a simple wrapper around a string that represents the format of the payload for a packet type in PDL
// source definitions and document generation. It provides a String method to retrieve the underlying string value.
type PayloadFormat struct {
	value string
}

// String returns the string representation of the PayloadFormat, which is simply the underlying string value.
func (p PayloadFormat) String() string {
	return p.value
}

var (
	DecimalFormat    PayloadFormat = PayloadFormat{value: token.DEC_FORMAT.String()}   // Decimal format for numeric values. e.g. 255
	HexFormat        PayloadFormat = PayloadFormat{value: token.HEX_FORMAT.String()}   // Hexadecimal format for numeric values. e.g. 0xFF
	BinaryFormat     PayloadFormat = PayloadFormat{value: token.BIN_FORMAT.String()}   // Binary format for numeric values. e.g. 0b1010
	BooleanFormat    PayloadFormat = PayloadFormat{value: token.BOOL_FORMAT.String()}  // Boolean format for true/false values. e.g. true
	ASCIIFormat      PayloadFormat = PayloadFormat{value: token.ASCII_FORMAT.String()} // ASCII format for text values. e.g. "Hello"
	UTF8Format       PayloadFormat = PayloadFormat{value: token.UTF8_FORMAT.String()}  // UTF-8 format for text values. e.g. "Hello"
	IPv4Format       PayloadFormat = PayloadFormat{value: token.IP4_FORMAT.String()}   // IPv4 format for IP addresses. e.g. 192.168.0.1
	IPv6Format       PayloadFormat = PayloadFormat{value: token.IP6_FORMAT.String()}   // IPv6 format for IP addresses. e.g. 2001:0db8:85a3:0000:0000:8a2e:0370:7334
	MACAddressFormat PayloadFormat = PayloadFormat{value: token.MAC_FORMAT.String()}   // MAC address format for network interfaces. e.g. 00:1A:2B:3C:4D:5E
	ArrayFormat      PayloadFormat = PayloadFormat{value: token.ARRAY_FORMAT.String()} // Array format for lists of values. e.g. [1, 2, 3]
)
