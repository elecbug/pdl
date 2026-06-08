package standard

import "github.com/elecbug/pdl"

func ARPPDL() pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.ARP.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    hardware_type        from 0   length 16
    protocol_type        from 16  length 16

    hardware_size        from 32  length 8
    protocol_size        from 40  length 8

    operation            from 48  length 16

    sender_hardware      from 64  length 48
    sender_protocol      from 112 length 32

    target_hardware      from 144 length 48
    target_protocol      from 192 length 32
}

out json {
    hardware_type arp.hardware_type {
        1       : "Ethernet"
        default : "Unknown"
    }

    protocol_type arp.protocol_type {
        0x0800  : "IPv4"
        default : "Unknown"
    }

    hardware_size arp.hardware_size DEC
    protocol_size arp.protocol_size DEC

    operation arp.operation {
        1       : "Request"
        2       : "Reply"
        default : "Unknown"
    }

    sender_hardware arp.sender.mac MAC
    sender_protocol arp.sender.ip  IP4

    target_hardware arp.target.mac MAC
    target_protocol arp.target.ip  IP4
}
`)
}
