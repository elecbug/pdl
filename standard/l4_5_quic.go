package standard

import "github.com/elecbug/pdl"

func QUICLongHeaderPDL(payload pdl.PayloadFormat) pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.QUICLong.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    first_byte      from 0  length 8

    header_form     from 0  length 1
    fixed_bit       from 1  length 1
    packet_type     from 2  length 2
    type_bits       from 4  length 4

    version         from 8  length 32

    dcid_len        from 40 length 8
    dcid            from 48 length (*dcid_len * 8)

    scid_len        from (48 + *dcid_len * 8) length 8
    scid            from (56 + *dcid_len * 8) length (*scid_len * 8)

    payload         from (56 + *dcid_len * 8 + *scid_len * 8) to end
}

out json {
    first_byte      quic.first_byte HEX

    header_form     quic.header_form BOOL
    fixed_bit       quic.fixed_bit   BOOL

    packet_type quic.packet_type {
        0       : "Initial"
        1       : "0-RTT"
        2       : "Handshake"
        3       : "Retry"
        default : "Unknown"
    }

    type_bits       quic.type_bits HEX

    version         quic.version HEX

    dcid_len        quic.destination_connection_id_length DEC
    dcid            quic.destination_connection_id HEX

    scid_len        quic.source_connection_id_length DEC
    scid            quic.source_connection_id HEX

    payload         quic.payload ` + payload.String() + `
}
`)
}

func QUICShortHeaderPDL(payload pdl.PayloadFormat) pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.QUICShort.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    first_byte        from 0 length 8

    header_form       from 0 length 1
    fixed_bit         from 1 length 1
    spin_bit          from 2 length 1
    reserved_bits     from 3 length 2
    key_phase         from 5 length 1
    packet_number_len from 6 length 2

    payload           from 8 to end
}

out json {
    first_byte        quic.first_byte HEX

    header_form       quic.header_form BOOL
    fixed_bit         quic.fixed_bit   BOOL
    spin_bit          quic.spin_bit    BOOL

    reserved_bits     quic.reserved_bits HEX
    key_phase         quic.key_phase BOOL

    packet_number_len quic.packet_number_length_raw DEC

    payload           quic.payload ` + payload.String() + `
}
`)
}
