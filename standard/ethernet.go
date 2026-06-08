package standard

import "github.com/elecbug/pdl"

func EthernetPDL() pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.Ethernet.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    dst_mac    from 0  length 48
    src_mac    from 48 length 48
    ether_type from 96 length 16

    payload    from 112 to end
}

out json {
    dst_mac    ethernet.destination MAC
    src_mac    ethernet.source      MAC

    ether_type ethernet.type {
        0x0800  : "IPv4"
        0x86DD  : "IPv6"
        0x0806  : "ARP"
        default : "Unknown"
    }

    payload ethernet.payload as switch *ether_type {
        0x0800  : IPv4
        0x86DD  : IPv6
        0x0806  : ARP
        default : HEX
    }
}
`)
}
