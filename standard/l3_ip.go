package standard

import "github.com/elecbug/pdl"

func IPv4Source() pdl.Source {
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
    version         ipv4.version              DEC
    ihl             ipv4.header_length_words  DEC
    dscp            ipv4.dscp                 DEC
    ecn             ipv4.ecn                  DEC

    total_length    ipv4.total_length         DEC
    identification  ipv4.identification       HEX

    fragment_offset ipv4.fragment_offset      DEC

    ttl             ipv4.ttl                  DEC

    protocol ipv4.protocol {
        1       : "ICMP"
        6       : "TCP"
        17      : "UDP"
        default : "Unknown"
    }

    checksum        ipv4.checksum             HEX

    src_ip          ipv4.source_ip            IP4
    dst_ip          ipv4.destination_ip       IP4

    options         ipv4.options              HEX

    payload ipv4.payload as switch *protocol {
        1       : ` + pdl.ICMP.String() + `
        6       : ` + pdl.TCP.String() + `
        17      : ` + pdl.UDP.String() + `
        default : ` + pdl.HexFormat.String() + `
    }

    flags<0> ipv4.flag.reserved {
        0 : false
        1 : true
    }

    flags<1> ipv4.flag.dont_fragment {
        0 : false
        1 : true
    }

    flags<2> ipv4.flag.more_fragments {
        0 : false
        1 : true
    }
}
`)
}

func IPv6Source() pdl.Source {
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
    version        ipv6.version          DEC
    traffic_class  ipv6.traffic_class    DEC
    flow_label     ipv6.flow_label       HEX

    payload_length ipv6.payload_length   DEC

    next_header ipv6.next_header {
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

    hop_limit      ipv6.hop_limit        DEC

    src_ip         ipv6.source_ip        IP6
    dst_ip         ipv6.destination_ip   IP6

    payload ipv6.payload as switch *next_header {
        6       : ` + pdl.TCP.String() + `
        17      : ` + pdl.UDP.String() + `
        44      : ` + pdl.IPv6Fragment.String() + `
        58      : ` + pdl.ICMPv6.String() + `
        default : ` + pdl.HexFormat.String() + `
    }
}
`)
}

func IPv6FragmentSource() pdl.Source {
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
