package standard

import (
	"github.com/elecbug/pdl"
)

// StandardSource generates a PDL source string for a given packet type and payload format.
// It returns a pdl.Source that can be used to create a PDL document for decoding packets of the
// specified type and extracting JSON output in the specified format.
func StandardSource(packet pdl.Packet, payload pdl.PayloadFormat) pdl.Source {
	switch packet {
	case pdl.Ethernet:
		return ethernetPDL(payload)
	case pdl.IPv4:
		return ipv4PDL(payload)
	case pdl.IPv6:
		return ipv6PDL(payload)
	case pdl.TCP:
		return tcpPDL(payload)
	default:
		return pdl.NewSource(`packet ` + packet.String())
	}
}
