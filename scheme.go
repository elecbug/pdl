package pdl

import (
	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/extractor/json_out"
	"github.com/elecbug/pdl/internal/parser"
)

// Scheme represents a compiled PDL document set that can be used to decode packets and extract JSON output according to the defined rules in the PDL sources.
// It contains a reference to the DocumentSet that was generated from the provided PDL sources and can be used to perform packet decoding and JSON extraction.
type Scheme struct {
	set *document.DocumentSet
}

// GenerateScheme takes a main packet name and a variable number of PDL source definitions, compiles them into a DocumentSet,
// and returns a Scheme that can be used to decode packets and extract JSON output according to the defined rules in the PDL sources.
// It returns an error if there is an issue with parsing the PDL sources or generating the DocumentSet.
func GenerateScheme(main Packet, sources ...Source) (*Scheme, error) {
	var strs []string
	for _, src := range sources {
		strs = append(strs, src.String())
	}

	doc, err := parser.ParseWithMultiSources(main.String(), strs...)
	if err != nil {
		return nil, err
	}

	return &Scheme{set: doc}, nil
}

// ExtractJSON takes a byte slice representing a packet, decodes it according to the rules defined in the Scheme's DocumentSet,
// and extracts JSON output based on the decoded packet data.
// It returns the extracted JSON output as an interface{} and an error if there is an issue with decoding the packet or extracting the JSON output.
func (p *Scheme) ExtractJSON(packet []byte) (any, error) {
	root := p.set.Root
	result, err := decoder.DecodeWithSet(p.set, root, packet)
	if err != nil {
		return nil, err
	}

	jsonRes, err := json_out.BuildJSONWithSet(p.set, root, result)
	if err != nil {
		return nil, err
	}

	return jsonRes, nil
}
