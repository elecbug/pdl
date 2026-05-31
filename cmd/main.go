package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/extractor"
	"github.com/elecbug/pdl/internal/parser"
)

func main() {
	var now time.Time

	srcFile := "doc/eg_tcp.pdl"
	binFile := "tmp/eg_tcp.bpdl"

	now = time.Now()
	src, err := os.ReadFile(srcFile)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", srcFile, err)
	}
	log.Printf("Reading file: %v", time.Since(now))

	var doc *document.Document
	{
		now = time.Now()
		doc, err = parser.ParseString(string(src))
		if err != nil {
			log.Fatalf("failed to parse PDL file %s: %v", srcFile, err)
		}
		log.Printf("Parsing: %v", time.Since(now))

		now = time.Now()
		jsonData, err := doc.Serialize()
		if err != nil {
			log.Fatalf("failed to serialize document: %v", err)
		}

		err = os.WriteFile(binFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("failed to write JSON file: %v", err)
		}
		log.Printf("Serialization: %v", time.Since(now))
	}
	{
		now = time.Now()
		jsonData, err := os.ReadFile(binFile)
		if err != nil {
			log.Fatalf("failed to read JSON file: %v", err)
		}

		doc, err = document.Deserialize(jsonData)
		if err != nil {
			log.Fatalf("failed to deserialize JSON file: %v", err)
		}
		log.Printf("Deserialization: %v", time.Since(now))
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

	now = time.Now()
	result, err := decoder.Decode(doc, packet)
	if err != nil {
		log.Fatalf("failed to decode packet: %v", err)
	}
	log.Printf("Decoding: %v", time.Since(now))

	now = time.Now()
	obj, err := extractor.BuildJSON(doc, result)
	if err != nil {
		log.Fatalf("failed to build JSON: %v", err)
	}
	log.Printf("Building JSON: %v", time.Since(now))

	now = time.Now()
	jsonData, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal JSON: %v", err)
	}
	log.Printf("Marshaling JSON: %v", time.Since(now))

	now = time.Now()
	err = os.WriteFile("tmp/log", jsonData, 0644)
	if err != nil {
		log.Fatalf("failed to write JSON file: %v", err)
	}
	log.Printf("Writing JSON file: %v", time.Since(now))
}
