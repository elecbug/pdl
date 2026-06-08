package pdl

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
	DecimalFormat    PayloadFormat = PayloadFormat{value: "DEC"}   // Decimal format for numeric values. e.g. 255
	HexFormat        PayloadFormat = PayloadFormat{value: "HEX"}   // Hexadecimal format for numeric values. e.g. 0xFF
	BinaryFormat     PayloadFormat = PayloadFormat{value: "BIN"}   // Binary format for numeric values. e.g. 0b1010
	BooleanFormat    PayloadFormat = PayloadFormat{value: "BOOL"}  // Boolean format for true/false values. e.g. true
	ASCIIFormat      PayloadFormat = PayloadFormat{value: "ASCII"} // ASCII format for text values. e.g. "Hello"
	UTF8Format       PayloadFormat = PayloadFormat{value: "UTF8"}  // UTF-8 format for text values. e.g. "Hello"
	IPv4Format       PayloadFormat = PayloadFormat{value: "IP4"}   // IPv4 format for IP addresses. e.g. 192.168.0.1
	IPv6Format       PayloadFormat = PayloadFormat{value: "IP6"}   // IPv6 format for IP addresses. e.g. 2001:0db8:85a3:0000:0000:8a2e:0370:7334
	MACAddressFormat PayloadFormat = PayloadFormat{value: "MAC"}   // MAC address format for network interfaces. e.g. 00:1A:2B:3C:4D:5E
)

var (
	EthernetPayloadFormat     PayloadFormat = PayloadFormat{value: "as " + Ethernet.String()}     // Payload format for Ethernet packet payloads
	IPv4PayloadFormat         PayloadFormat = PayloadFormat{value: "as " + IPv4.String()}         // Payload format for IPv4 packet payloads
	IPv6PayloadFormat         PayloadFormat = PayloadFormat{value: "as " + IPv6.String()}         // Payload format for IPv6 packet payloads
	IPv6FragmentPayloadFormat PayloadFormat = PayloadFormat{value: "as " + IPv6Fragment.String()} // Payload format for IPv6 fragment packet payloads
	TCPPayloadFormat          PayloadFormat = PayloadFormat{value: "as " + TCP.String()}          // Payload format for TCP packet payloads
	UDPPayloadFormat          PayloadFormat = PayloadFormat{value: "as " + UDP.String()}          // Payload format for UDP packet payloads
	QUICLongPayloadFormat     PayloadFormat = PayloadFormat{value: "as " + QUICLong.String()}     // Payload format for QUIC long header packet payloads
	QUICShortPayloadFormat    PayloadFormat = PayloadFormat{value: "as " + QUICShort.String()}    // Payload format for QUIC short header packet payloads
)
