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
    data_offset    from 96  length 4
    options        from 160 length (*data_offset * 32 - 160)
    payload        from (*data_offset * 32) to end
}

out json {
    src_port source_port DEC
    dst_port destination_port DEC
}
`

	doc, err := pdl.ParseString(src)
	if err != nil {
		log.Fatal(err)
	}

	// TCP header:
	// src=80, dst=443, data_offset=5
	packet := []byte{
		0x00, 0x50, // src_port = 80
		0x01, 0xbb, // dst_port = 443
		0x00, 0x00, 0x00, 0x00, // seq
		0x00, 0x00, 0x00, 0x00, // ack
		0x50, 0x02, // data_offset=5, flags=SYN
		0x20, 0x00, // window
		0x00, 0x00, // checksum
		0x00, 0x00, // urgent
		0xde, 0xad, 0xbe, 0xef, // payload
	}

	result, err := pdl.Decode(doc, packet)
	if err != nil {
		log.Fatal(err)
	}

	for name, v := range result.Values {
		fmt.Printf("%s = %d, bits=%X, len=%d\n", name, v.UInt, v.Bits, v.Len)
	}

	jsonObj, err := pdl.BuildJSON(
		doc,
		result,
	)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(
		jsonObj,
		"",
		"  ",
	)

	fmt.Println(string(b))
}
