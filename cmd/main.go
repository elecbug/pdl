package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"os"

	"github.com/elecbug/pdl"
	"github.com/elecbug/pdl/standard"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <hex_string>", os.Args[0])
	}

	pdl, err := pdl.Generate(
		standard.Ethernet,
		standard.Ethernet_PDL(standard.IPv4.WithAs()),
		standard.IPv4_PDL(standard.TCP.WithAs()),
		standard.TCP_PDL("HEX"),
	)
	if err != nil {
		log.Fatalf("failed to parse PDL file: %v", err)
	}

	packet, err := hex.DecodeString(os.Args[1])
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}

	result, err := pdl.ExtractJSON(packet)
	if err != nil {
		log.Fatalf("failed to extract JSON: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("tmp/tcp_log", jsonData, 0644)
	if err != nil {
		log.Fatalf("failed to write JSON file: %v", err)
	}
}
