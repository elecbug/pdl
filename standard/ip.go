package standard

import "github.com/elecbug/pdl"

func IPv4PDL(payload pdl.Payload) pdl.Source {
	return pdl.Source(`
packet ` + IPv4.String() + `

set mode BIG_ENDIAN MSB_FIRST

var {
    fixed_header_bits = 160
}

def {
    version        from 0   length 4
    ihl            from 4   length 4
    dscp           from 8   length 6
    ecn            from 14  length 2

    total_length   from 16  length 16
    identification from 32  length 16

    flags          from 48  length 3
    fragment_offset from 51 length 13

    ttl            from 64  length 8
    protocol       from 72  length 8
    checksum       from 80  length 16

    src_ip         from 96  length 32
    dst_ip         from 128 length 32

    options        from fixed_header_bits length (*ihl * 32 - fixed_header_bits)
    payload        from (*ihl * 32) to end
}

out json {
    version         version              DEC
    ihl             header.length_words  DEC
    dscp            dscp                 DEC
    ecn             ecn                  DEC

    total_length    total_length         DEC
    identification  identification       HEX

    fragment_offset fragment_offset      DEC

    ttl             ttl                  DEC

    protocol protocol {
        1  : "ICMP"
        6  : "TCP"
        17 : "UDP"
    }

    checksum        checksum             HEX

    src_ip          source_ip            IP4
    dst_ip          destination_ip       IP4

    options         options              HEX
    payload         payload              ` + string(payload) + `

    flags<0> flag.reserved {
        0 : false
        1 : true
    }

    flags<1> flag.dont_fragment {
        0 : false
        1 : true
    }

    flags<2> flag.more_fragments {
        0 : false
        1 : true
    }
}
`)
}
