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

def {
    src_port       from 0   length 16
    dst_port       from 16  length 16
    seq            from 32  length 32
    ack            from 64  length 32
    data_offset    from 96  length 4
    flags          from 104 length 8
    payload        from (*data_offset * 32) to end
}

out json {
    src_port source.port DEC
    dst_port destination.port DEC
    seq sequence.number DEC
    ack acknowledgment.number DEC

    flags<6> flags.syn {
        0 : false
        1 : true
    }

    payload payload.raw HEX
}
`

	doc, err := pdl.ParseString(src)
	if err != nil {
		log.Fatal(err)
	}

	packet := []byte{
		0x00, 0x50,
		0x01, 0xbb,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02,
		0x50, 0x02,
		0x20, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0xde, 0xad, 0xbe, 0xef,
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