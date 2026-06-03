package pdl

import (
	"fmt"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/extractor/json_out"
	"github.com/elecbug/pdl/internal/parser"
)

// PDL represents a parsed document set that defines how to decode packets and extract JSON output based on the defined packet types and payloads.
// It embeds the DocumentSet structure, allowing it to hold multiple documents and a designated root document for decoding.
type PDL struct {
	set *document.DocumentSet
}

// PDLSource is a type alias for string, representing the source of a PDLSource document, which can be used to generate a Document structure.
type PDLSource string

// Generate takes a main PacketType and one or more PDLSource strings, parses them into a structured DocumentSet,
// and returns a pointer to a PDL representing the parsed document set. It returns an error if parsing fails.
func Generate(main PacketType, sources ...PDLSource) (*PDL, error) {
	var strSrcs []string
	for _, src := range sources {
		strSrcs = append(strSrcs, string(src))
	}

	doc, err := parser.ParseWithMultiSources(main.String(), strSrcs...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	return &PDL{set: doc}, nil
}

// ExtractJSON takes a byte slice representing a packet, decodes it according to the PDL document's rules,
// and constructs a JSON-compatible Go value based on the document's output configuration. It returns the
// resulting JSON value or an error if decoding or JSON construction fails.
func (doc *PDL) ExtractJSON(packet []byte) (any, error) {
	root := doc.set.Root
	result, err := decoder.Decode(root, packet)
	if err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	jsonRes, err := json_out.BuildJSONWithSet(doc.set, root, result)
	if err != nil {
		return nil, fmt.Errorf("failed to build JSON with set: %w", err)
	}

	return jsonRes, nil
}

// Payload is a type alias for string, representing the name of the main packet type defined in the PDL document.
type Payload string

// String returns the string representation of the Payload, which is simply the underlying string value.
func (p Payload) String() string {
	return string(p)
}

// PacketType is a type alias for string, representing the name of a packet type that can be used in PDL source definitions and document generation.
type PacketType string

// AsPayload returns a Payload that represents the packet type as a payload, which can be used in PDL source
// definitions to specify the main packet type for decoding.
func (p PacketType) AsPayload() Payload {
	return Payload("as " + p.String())
}

// String returns the string representation of the PacketType, which is simply the underlying string value.
func (p PacketType) String() string {
	return string(p)
}
