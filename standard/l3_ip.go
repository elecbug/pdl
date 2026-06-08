package standard

import "github.com/elecbug/pdl"

func IPv4PDL() pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.IPv4.String() + `

set mode BIG_ENDIAN MSB_FIRST

var {
    fixed_header_bits = 160
}

def {
    version         from 0   length 4
    ihl             from 4   length 4
    dscp            from 8   length 6
    ecn             from 14  length 2

    total_length    from 16  length 16
    identification  from 32  length 16

    flags           from 48  length 3
    fragment_offset from 51  length 13

    ttl             from 64  length 8
    protocol        from 72  length 8
    checksum        from 80  length 16

    src_ip          from 96  length 32
    dst_ip          from 128 length 32

    options         from fixed_header_bits length (*ihl * 32 - fixed_header_bits)
    payload         from (*ihl * 32) to end
}

out json {
    version         ip.version              DEC
    ihl             ip.header_length_words  DEC
    dscp            ip.dscp                 DEC
    ecn             ip.ecn                  DEC

    total_length    ip.total_length         DEC
    identification  ip.identification       HEX

    fragment_offset ip.fragment_offset      DEC

    ttl             ip.ttl                  DEC

    protocol ip.protocol {
        1       : "ICMP"
        6       : "TCP"
        17      : "UDP"
        default : "Unknown"
    }

    checksum        ip.checksum             HEX

    src_ip          ip.source_ip            IP4
    dst_ip          ip.destination_ip       IP4

    options         ip.options              HEX

    payload ip.payload as switch *protocol {
        1       : ` + pdl.ICMP.String() + `
        6       : ` + pdl.TCP.String() + `
        17      : ` + pdl.UDP.String() + `
        default : ` + pdl.HexFormat.String() + `
    }

    flags<0> ip.flag.reserved {
        0 : false
        1 : true
    }

    flags<1> ip.flag.dont_fragment {
        0 : false
        1 : true
    }

    flags<2> ip.flag.more_fragments {
        0 : false
        1 : true
    }
}
`)
}

func IPv6PDL() pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.IPv6.String() + `

set mode BIG_ENDIAN MSB_FIRST

var {
    fixed_header_bits = 320
}

def {
    version        from 0   length 4
    traffic_class  from 4   length 8
    flow_label     from 12  length 20

    payload_length from 32  length 16
    next_header    from 48  length 8
    hop_limit      from 56  length 8

    src_ip         from 64  length 128
    dst_ip         from 192 length 128

    payload        from fixed_header_bits to end
}

out json {
    version        ip.version          DEC
    traffic_class  ip.traffic_class    DEC
    flow_label     ip.flow_label       HEX

    payload_length ip.payload_length   DEC

    next_header ip.next_header {
        0       : "Hop-by-Hop Options"
        1       : "ICMP"
        6       : "TCP"
        17      : "UDP"
        41      : "IPv6"
        43      : "Routing"
        44      : "Fragment"
        50      : "ESP"
        51      : "AH"
        58      : "ICMPv6"
        59      : "No Next Header"
        60      : "Destination Options"
        132     : "SCTP"
        default : "Unknown"
    }

    hop_limit      ip.hop_limit        DEC

    src_ip         ip.source_ip        IP6
    dst_ip         ip.destination_ip   IP6

    payload ip.payload as switch *next_header {
        6       : ` + pdl.TCP.String() + `
        17      : ` + pdl.UDP.String() + `
        44      : ` + pdl.IPv6Fragment.String() + `
        58      : ` + pdl.ICMPv6.String() + `
        default : ` + pdl.HexFormat.String() + `
    }
}
`)
}

func IPv6FragmentPDL() pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.IPv6Fragment.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    next_header       from 0  length 8
    reserved_octet    from 8  length 8

    fragment_offset   from 16 length 13
    reserved_bits     from 29 length 2
    more_fragments    from 31 length 1

    identification    from 32 length 32

    payload           from 64 to end
}

out json {
    next_header fragment.next_header {
        6       : "TCP"
        17      : "UDP"
        44      : "Fragment"
        58      : "ICMPv6"
        default : "Unknown"
    }

    reserved_octet    fragment.reserved_octet HEX
    fragment_offset   fragment.offset         DEC
    reserved_bits     fragment.reserved_bits  HEX
    more_fragments    fragment.more_fragments BOOL
    identification    fragment.identification HEX

    payload fragment.payload as switch *next_header {
        6       : ` + pdl.TCP.String() + `
        17      : ` + pdl.UDP.String() + `
        44      : ` + pdl.IPv6Fragment.String() + `
        58      : ` + pdl.ICMPv6.String() + `
        default : ` + pdl.HexFormat.String() + `
    }
}
`)
}
