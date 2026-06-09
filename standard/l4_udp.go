package standard

import "github.com/elecbug/pdl"

func UDPSource(payload pdl.PayloadFormat) pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.UDP.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    src_port    from 0  length 16
    dst_port    from 16 length 16
    len         from 32 length 16
    checksum    from 48 length 16

    payload     from 64 to end
}

out json {
    src_port    udp.source_port      DEC
    dst_port    udp.destination_port DEC
    len         udp.length           DEC
    checksum    udp.checksum         HEX

    payload     udp.payload          ` + payload.String() + `
}
`)
}
