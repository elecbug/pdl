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
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s <out_file> <hex_string>", os.Args[0])
	}

	scheme, err := pdl.GenerateScheme(
		pdl.Ethernet,
		standard.CommonSources()...,
	)
	if err != nil {
		log.Fatalf("failed to parse PDL file: %v", err)
	}

	outFile := os.Args[1]

	packet, err := hex.DecodeString(os.Args[2])
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}

	result, err := scheme.ExtractJSON(packet)
	if err != nil {
		log.Fatalf("failed to extract JSON: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(outFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("failed to write JSON to file: %v", err)
	}
}
