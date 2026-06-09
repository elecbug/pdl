package standard

import "github.com/elecbug/pdl"

func ICMPSource(paylaod pdl.PayloadFormat) pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.ICMP.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    type      from 0  length 8
    code      from 8  length 8
    checksum  from 16 length 16

    rest      from 32 length 32

    payload   from 64 to end
}

out json {
    type icmp.type {
        0  : "Echo Reply"
        3  : "Destination Unreachable"
        5  : "Redirect"
        8  : "Echo Request"
        11 : "Time Exceeded"
        12 : "Parameter Problem"
        default : "Unknown"
    }

    code      icmp.code      DEC
    checksum  icmp.checksum  HEX
    rest      icmp.rest      HEX

    payload   icmp.payload   ` + paylaod.String() + `
}
`)
}

func ICMPv6Source(payload pdl.PayloadFormat) pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.ICMPv6.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    type      from 0  length 8
    code      from 8  length 8
    checksum  from 16 length 16

    rest      from 32 length 32

    payload   from 64 to end
}

out json {
    type icmpv6.type {
        1   : "Destination Unreachable"
        2   : "Packet Too Big"
        3   : "Time Exceeded"
        4   : "Parameter Problem"

        128 : "Echo Request"
        129 : "Echo Reply"

        133 : "Router Solicitation"
        134 : "Router Advertisement"

        135 : "Neighbor Solicitation"
        136 : "Neighbor Advertisement"

        137 : "Redirect"

        default : "Unknown"
    }

    code      icmpv6.code      DEC
    checksum  icmpv6.checksum  HEX
    rest      icmpv6.rest      HEX

    payload   icmpv6.payload   ` + payload.String() + `
}
`)
}
