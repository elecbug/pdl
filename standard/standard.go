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
		UDPSource(),
		ICMPSource(pdl.HexFormat),
		ICMPv6Source(pdl.HexFormat),
		DNSSource(pdl.HexFormat),
		QUICSource(),
		QUICLongHeaderSource(pdl.HexFormat),
		QUICInitialLikeSource(pdl.HexFormat),
		QUICRetrySource(pdl.HexFormat),
		QUICShortHeaderSource(pdl.HexFormat),
	}
}
