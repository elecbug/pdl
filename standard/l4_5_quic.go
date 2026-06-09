package standard

import (
	"fmt"
	"strings"

	"github.com/elecbug/pdl"
)

func QUICSource() pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.QUIC.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    header_form from 0 length 1
    payload     from 0 to end
}

out json {
    payload quic as switch *header_form {
        0       : ` + pdl.QUICShort.String() + `
        1       : ` + pdl.QUICLong.String() + `
        default : ` + pdl.HexFormat.String() + `
    }
}
`)
}

func QUICLongHeaderSource(payload pdl.PayloadFormat) pdl.Source {
	_ = payload

	return pdl.NewSource(`
packet ` + pdl.QUICLong.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    first_byte         from 0 length 8

    header_form        from 0 length 1
    fixed_bit          from 1 length 1
    packet_type        from 2 length 2
    reserved_bits      from 4 length 2
    packet_number_len  from 6 length 2

    payload            from 0 to end
}

out json {
    payload long_header as switch *packet_type {
        0       : ` + pdl.QUICInitialLike.String() + `
        1       : ` + pdl.QUICInitialLike.String() + `
        2       : ` + pdl.QUICInitialLike.String() + `
        3       : ` + pdl.QUICRetry.String() + `
        default : ` + pdl.HexFormat.String() + `
    }
}
`)
}

func QUICInitialLikeSource(payload pdl.PayloadFormat) pdl.Source {
	packetNumberSwitch := buildPacketNumberSwitch()
	payloadSwitch := buildPayloadSwitch()

	return pdl.NewSource(`
packet ` + pdl.QUICInitialLike.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    first_byte         from 0 length 8

    header_form        from 0 length 1
    fixed_bit          from 1 length 1
    packet_type        from 2 length 2
    reserved_bits      from 4 length 2
    packet_number_len  from 6 length 2

    version            from 8 length 32

    dcid_len           from 40 length 8
    dcid               from 48 length (*dcid_len * 8)

    scid_len           from (48 + *dcid_len * 8) length 8
    scid               from (56 + *dcid_len * 8) length (*scid_len * 8)

    token_len_prefix   from (56 + *dcid_len * 8 + *scid_len * 8) length 2

    token_len switch *token_len_prefix {
        0 : from (56 + *dcid_len * 8 + *scid_len * 8 + 2) length 6
        1 : from (56 + *dcid_len * 8 + *scid_len * 8 + 2) length 14
        2 : from (56 + *dcid_len * 8 + *scid_len * 8 + 2) length 30
        3 : from (56 + *dcid_len * 8 + *scid_len * 8 + 2) length 62
    }

    token switch *token_len_prefix {
        0 : from (56 + *dcid_len * 8 + *scid_len * 8 + 8)  length (*token_len * 8)
        1 : from (56 + *dcid_len * 8 + *scid_len * 8 + 16) length (*token_len * 8)
        2 : from (56 + *dcid_len * 8 + *scid_len * 8 + 32) length (*token_len * 8)
        3 : from (56 + *dcid_len * 8 + *scid_len * 8 + 64) length (*token_len * 8)
    }

    length_prefix switch *token_len_prefix {
        0 : from (56 + *dcid_len * 8 + *scid_len * 8 + 8  + *token_len * 8) length 2
        1 : from (56 + *dcid_len * 8 + *scid_len * 8 + 16 + *token_len * 8) length 2
        2 : from (56 + *dcid_len * 8 + *scid_len * 8 + 32 + *token_len * 8) length 2
        3 : from (56 + *dcid_len * 8 + *scid_len * 8 + 64 + *token_len * 8) length 2
    }

    len switch *length_prefix {
        val == 0 && *token_len_prefix == 0 : from (56 + *dcid_len * 8 + *scid_len * 8 + 8  + *token_len * 8 + 2) length 6
        val == 1 && *token_len_prefix == 0 : from (56 + *dcid_len * 8 + *scid_len * 8 + 8  + *token_len * 8 + 2) length 14
        val == 2 && *token_len_prefix == 0 : from (56 + *dcid_len * 8 + *scid_len * 8 + 8  + *token_len * 8 + 2) length 30
        val == 3 && *token_len_prefix == 0 : from (56 + *dcid_len * 8 + *scid_len * 8 + 8  + *token_len * 8 + 2) length 62

        val == 0 && *token_len_prefix == 1 : from (56 + *dcid_len * 8 + *scid_len * 8 + 16 + *token_len * 8 + 2) length 6
        val == 1 && *token_len_prefix == 1 : from (56 + *dcid_len * 8 + *scid_len * 8 + 16 + *token_len * 8 + 2) length 14
        val == 2 && *token_len_prefix == 1 : from (56 + *dcid_len * 8 + *scid_len * 8 + 16 + *token_len * 8 + 2) length 30
        val == 3 && *token_len_prefix == 1 : from (56 + *dcid_len * 8 + *scid_len * 8 + 16 + *token_len * 8 + 2) length 62

        val == 0 && *token_len_prefix == 2 : from (56 + *dcid_len * 8 + *scid_len * 8 + 32 + *token_len * 8 + 2) length 6
        val == 1 && *token_len_prefix == 2 : from (56 + *dcid_len * 8 + *scid_len * 8 + 32 + *token_len * 8 + 2) length 14
        val == 2 && *token_len_prefix == 2 : from (56 + *dcid_len * 8 + *scid_len * 8 + 32 + *token_len * 8 + 2) length 30
        val == 3 && *token_len_prefix == 2 : from (56 + *dcid_len * 8 + *scid_len * 8 + 32 + *token_len * 8 + 2) length 62

        val == 0 && *token_len_prefix == 3 : from (56 + *dcid_len * 8 + *scid_len * 8 + 64 + *token_len * 8 + 2) length 6
        val == 1 && *token_len_prefix == 3 : from (56 + *dcid_len * 8 + *scid_len * 8 + 64 + *token_len * 8 + 2) length 14
        val == 2 && *token_len_prefix == 3 : from (56 + *dcid_len * 8 + *scid_len * 8 + 64 + *token_len * 8 + 2) length 30
        val == 3 && *token_len_prefix == 3 : from (56 + *dcid_len * 8 + *scid_len * 8 + 64 + *token_len * 8 + 2) length 62
    }

` + packetNumberSwitch + `

` + payloadSwitch + `
}

out json {
    first_byte        initial.first_byte HEX

    header_form       initial.header_form BOOL
    fixed_bit         initial.fixed_bit   BOOL

    packet_type initial.packet_type {
        0       : "Initial"
        1       : "0-RTT"
        2       : "Handshake"
        default : "Unknown"
    }

    reserved_bits     initial.reserved_bits HEX
    packet_number_len initial.packet_number_length_raw DEC

    version           initial.version HEX

    dcid_len          initial.destination_connection_id_length DEC
    dcid              initial.destination_connection_id HEX

    scid_len          initial.source_connection_id_length DEC
    scid              initial.source_connection_id HEX

    token_len_prefix  initial.token_length_prefix DEC
    token_len         initial.token_length DEC
    token             initial.token HEX

    length_prefix     initial.length_prefix DEC
    len               initial.length DEC

    packet_number     initial.packet_number DEC

    payload           initial.payload ` + payload.String() + `
}
`)
}

func QUICRetrySource(payload pdl.PayloadFormat) pdl.Source {
	_ = payload

	return pdl.NewSource(`
packet ` + pdl.QUICRetry.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    first_byte         from 0 length 8

    header_form        from 0 length 1
    fixed_bit          from 1 length 1
    packet_type        from 2 length 2
    type_bits          from 4 length 4

    version            from 8 length 32

    dcid_len           from 40 length 8
    dcid               from 48 length (*dcid_len * 8)

    scid_len           from (48 + *dcid_len * 8) length 8
    scid               from (56 + *dcid_len * 8) length (*scid_len * 8)

    token_start        from (56 + *dcid_len * 8 + *scid_len * 8) length 0

    retry_token         from *token_start to (end - 128)
    retry_integrity_tag from (end - 127) to end
}

out json {
    first_byte         retry.first_byte HEX

    header_form        retry.header_form BOOL
    fixed_bit          retry.fixed_bit   BOOL

    packet_type retry.packet_type {
        3       : "Retry"
        default : "Unknown"
    }

    type_bits          retry.type_bits HEX

    version            retry.version HEX

    dcid_len           retry.destination_connection_id_length DEC
    dcid               retry.destination_connection_id HEX

    scid_len           retry.source_connection_id_length DEC
    scid               retry.source_connection_id HEX

    retry_token         retry.retry_token HEX
    retry_integrity_tag retry.retry_integrity_tag HEX
}
`)
}

func QUICShortHeaderSource(payload pdl.PayloadFormat) pdl.Source {
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

    packet_number switch *packet_number_len {
        0 : from 8 length 8
        1 : from 8 length 16
        2 : from 8 length 24
        3 : from 8 length 32
    }

    payload switch *packet_number_len {
        0 : from 16 to end
        1 : from 24 to end
        2 : from 32 to end
        3 : from 40 to end
    }
}

out json {
    first_byte        short_header.first_byte HEX

    header_form       short_header.header_form BOOL
    fixed_bit         short_header.fixed_bit   BOOL
    spin_bit          short_header.spin_bit    BOOL

    reserved_bits     short_header.reserved_bits HEX
    key_phase         short_header.key_phase BOOL

    packet_number_len short_header.packet_number_length_raw DEC
    packet_number     short_header.packet_number DEC

    payload           short_header.payload ` + payload.String() + `
}
`)
}

func buildPacketNumberSwitch() string {
	var sb strings.Builder

	sb.WriteString(`
    packet_number switch *packet_number_len {
`)

	tokenLenBits := []int{8, 16, 32, 64}
	lengthBits := []int{8, 16, 32, 64}
	pnBits := []int{8, 16, 24, 32}

	for tpIdx, tokenLenBit := range tokenLenBits {
		for lpIdx, lengthBit := range lengthBits {
			for pnIdx, pnBit := range pnBits {
				sb.WriteString(fmt.Sprintf(
					`        val == %d && *token_len_prefix == %d && *length_prefix == %d : from (56 + *dcid_len * 8 + *scid_len * 8 + %d + *token_len * 8 + %d) length %d
`,
					pnIdx,
					tpIdx,
					lpIdx,
					tokenLenBit,
					lengthBit,
					pnBit,
				))
			}
		}
	}

	sb.WriteString("    }\n")
	return sb.String()
}

func buildPayloadSwitch() string {
	var sb strings.Builder

	sb.WriteString(`
    payload switch *packet_number_len {
`)

	tokenLenBits := []int{8, 16, 32, 64}
	lengthBits := []int{8, 16, 32, 64}
	pnBits := []int{8, 16, 24, 32}

	for tpIdx, tokenLenBit := range tokenLenBits {
		for lpIdx, lengthBit := range lengthBits {
			for pnIdx, pnBit := range pnBits {
				sb.WriteString(fmt.Sprintf(
					`        val == %d && *token_len_prefix == %d && *length_prefix == %d : from (56 + *dcid_len * 8 + *scid_len * 8 + %d + *token_len * 8 + %d + %d) to end
`,
					pnIdx,
					tpIdx,
					lpIdx,
					tokenLenBit,
					lengthBit,
					pnBit,
				))
			}
		}
	}

	sb.WriteString("    }\n")
	return sb.String()
}
