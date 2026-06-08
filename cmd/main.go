package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elecbug/pdl"
	"github.com/elecbug/pdl/standard"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <hex_string>", os.Args[0])
	}

	scheme, err := pdl.GenerateScheme(
		pdl.Ethernet,
		standard.EthernetPDL(),
		standard.ARPPDL(),
		standard.IPv4PDL(),
		standard.IPv6PDL(),
		standard.IPv6FragmentPDL(),
		standard.TCPPDL(pdl.HexFormat),
		standard.UDPPDL(pdl.HexFormat),
	)
	if err != nil {
		log.Fatalf("failed to parse PDL file: %v", err)
	}

	packet, err := hex.DecodeString(os.Args[1])
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

	fmt.Println(string(jsonData))
}
