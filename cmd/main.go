package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elecbug/pdl/internal/pdl"
)

func main() {
	srcFile := "doc/eg_tcp.pdl"
	src, err := os.ReadFile(srcFile)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", srcFile, err)
	}

	doc, err := pdl.ParseString(string(src))
	if err != nil {
		log.Fatalf("failed to parse PDL file %s: %v", srcFile, err)
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
		log.Fatalf("failed to decode packet: %v", err)
	}

	obj, err := pdl.BuildJSON(doc, result)
	if err != nil {
		log.Fatalf("failed to build JSON: %v", err)
	}

	b, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(b))
}
