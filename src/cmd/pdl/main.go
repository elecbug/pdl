package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elecbug/pdl/src/internal/pdl"
)

func main() {
	src := `
packet TCP

set mode BIG_ENDIAN

var {
    fixed_header_bits = 160
}

def {
    src_port       from 0   length 16
    dst_port       from 16  length 16

    seq            from 32  length 32
    ack            from 64  length 32

    data_offset    from 96  length 4
    reserved       from 100 length 3
    ns             from 103 length 1
    flags          from 104 length 8

    window         from 112 length 16
    checksum       from 128 length 16
    urgent_pointer from 144 length 16

    options        from fixed_header_bits length (*data_offset * 32 - fixed_header_bits)
    payload        from (*data_offset * 32) to end
}

out json {
    src_port       source.port              DEC
    dst_port       destination.port         DEC

    seq            sequence_number          DEC
    ack            acknowledgment_number    DEC

    data_offset    header.length_words      DEC
    reserved       header.reserved          BIN
    ns             header.ns                BOOL

    window         window_size              DEC
    checksum       checksum                 HEX
    urgent_pointer urgent_pointer           DEC

    options        options                  HEX
    payload        payload                  HEX

    flags<0> flags.cwr {
        0 : false
        1 : true
    }

    flags<1> flags.ece {
        0 : false
        1 : true
    }

    flags<2> flags.urg {
        0 : false
        1 : true
    }

    flags<3> flags.ack {
        0 : false
        1 : true
    }

    flags<4> flags.psh {
        0 : false
        1 : true
    }

    flags<5> flags.rst {
        0 : false
        1 : true
    }

    flags<6> flags.syn {
        0 : false
        1 : true
    }

    flags<7> flags.fin {
        0 : false
        1 : true
    }
}
`

	doc, err := pdl.ParseString(src)
	if err != nil {
		log.Fatal(err)
	}

	packet := []byte{
		0x00, 0x50, // src_port = 80
		0x01, 0xbb, // dst_port = 443

		0x00, 0x00, 0x00, 0x01, // seq = 1
		0x00, 0x00, 0x00, 0x00, // ack = 0

		0x50, // data_offset=5, reserved=0, ns=0
		0x02, // flags=SYN

		0x20, 0x00, // window = 8192
		0xab, 0xcd, // checksum
		0x00, 0x00, // urgent pointer

		0xde, 0xad, 0xbe, 0xef, // payload
	}

	result, err := pdl.Decode(doc, packet)
	if err != nil {
		log.Fatal(err)
	}

	obj, err := pdl.BuildJSON(doc, result)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(b))
}
