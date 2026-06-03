package standard

import "github.com/elecbug/pdl"

const (
	Ethernet pdl.PacketType = "Ethernet"
	IPv4     pdl.PacketType = "IPv4"
	TCP      pdl.PacketType = "TCP"

	DecimalFormat    pdl.Payload = "DEC"
	HexFormat        pdl.Payload = "HEX"
	BinaryFormat     pdl.Payload = "BIN"
	BooleanFormat    pdl.Payload = "BOOL"
	ASCIIFormat      pdl.Payload = "ASCII"
	UTF8Format       pdl.Payload = "UTF8"
	IPv4Format       pdl.Payload = "IP4"
	IPv6Format       pdl.Payload = "IP6"
	MACAddressFormat pdl.Payload = "MAC"
)
