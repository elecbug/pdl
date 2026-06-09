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
    payload quic as switch *packet_type {
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
    first_byte        quic.first_byte HEX

    header_form       quic.header_form BOOL
    fixed_bit         quic.fixed_bit   BOOL

    packet_type quic.packet_type {
        0       : "Initial"
        1       : "0-RTT"
        2       : "Handshake"
        default : "Unknown"
    }

    reserved_bits     quic.reserved_bits HEX
    packet_number_len quic.packet_number_length_raw DEC

    version           quic.version HEX

    dcid_len          quic.destination_connection_id_length DEC
    dcid              quic.destination_connection_id HEX

    scid_len          quic.source_connection_id_length DEC
    scid              quic.source_connection_id HEX

    token_len_prefix  quic.token_length_prefix DEC
    token_len         quic.token_length DEC
    token             quic.token HEX

    length_prefix     quic.length_prefix DEC
    len               quic.length DEC

    packet_number     quic.packet_number DEC

    payload           quic.payload ` + payload.String() + `
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
    first_byte         quic.first_byte HEX

    header_form        quic.header_form BOOL
    fixed_bit          quic.fixed_bit   BOOL

    packet_type quic.packet_type {
        3       : "Retry"
        default : "Unknown"
    }

    type_bits          quic.type_bits HEX

    version            quic.version HEX

    dcid_len           quic.destination_connection_id_length DEC
    dcid               quic.destination_connection_id HEX

    scid_len           quic.source_connection_id_length DEC
    scid               quic.source_connection_id HEX

    retry_token         quic.retry_token HEX
    retry_integrity_tag quic.retry_integrity_tag HEX
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
    first_byte        quic.first_byte HEX

    header_form       quic.header_form BOOL
    fixed_bit         quic.fixed_bit   BOOL
    spin_bit          quic.spin_bit    BOOL

    reserved_bits     quic.reserved_bits HEX
    key_phase         quic.key_phase BOOL

    packet_number_len quic.packet_number_length_raw DEC
    packet_number     quic.packet_number DEC

    payload           quic.payload ` + payload.String() + `
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
