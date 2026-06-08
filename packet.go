package pdl

// Packet represents the name of a packet type that can be used in PDL source definitions and document generation.
// It provides a String method to retrieve the underlying string value and an AsPayload method to represent the packet type as a payload in PDL source definitions.
type Packet struct {
	value string
}

// NewPacket creates a new Packet instance from a given string, which represents the name of a packet type that
// can be used in PDL source definitions and document generation.
func NewPacket(value string) Packet {
	return Packet{value: value}
}

// String returns the string representation of the Packet, which is simply the underlying string value.
func (p Packet) String() string {
	return p.value
}

var (
	Ethernet     Packet = Packet{value: "Ethernet"}
	IPv4         Packet = Packet{value: "IPv4"}
	IPv6         Packet = Packet{value: "IPv6"}
	IPv6Fragment Packet = Packet{value: "IPv6Fragment"}
	TCP          Packet = Packet{value: "TCP"}
	UDP          Packet = Packet{value: "UDP"}
	QUICLong     Packet = Packet{value: "QUICLong"}
	QUICShort    Packet = Packet{value: "QUICShort"}
)
