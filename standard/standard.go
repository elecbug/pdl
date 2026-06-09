package standard

import "github.com/elecbug/pdl"

// CommonSources returns a list of PDL sources for common network protocols.
func CommonSources() []pdl.Source {
	return []pdl.Source{
		EthernetSource(),
		ARPSource(),
		IPv4Source(),
		IPv6Source(),
		IPv6FragmentSource(),
		TCPSource(pdl.HexFormat),
		UDPSource(pdl.HexFormat),
		ICMPSource(pdl.HexFormat),
		ICMPv6Source(pdl.HexFormat),
	}
}
