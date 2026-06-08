package standard

import "github.com/elecbug/pdl"

// GenericSources returns a list of PDL sources for common network protocols.
func GenericSources() []pdl.Source {
	return []pdl.Source{
		EthernetPDL(),
		ARPPDL(),
		IPv4PDL(),
		IPv6PDL(),
		IPv6FragmentPDL(),
		TCPPDL(pdl.HexFormat),
		UDPPDL(pdl.HexFormat),
		ICMPPDL(pdl.HexFormat),
		ICMPv6PDL(pdl.HexFormat),
	}
}
