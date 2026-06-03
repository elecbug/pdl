package pdl

import (
	"fmt"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/extractor/json_out"
	"github.com/elecbug/pdl/internal/parser"
)

// PDL represents a parsed PDL document set, which contains one or more packet definitions and their associated decoding rules.
type PDL document.DocumentSet

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

	return (*PDL)(doc), nil
}

// ExtractJSON takes a byte slice representing a packet, decodes it according to the PDL document's rules,
// and constructs a JSON-compatible Go value based on the document's output configuration. It returns the
// resulting JSON value or an error if decoding or JSON construction fails.
func (doc *PDL) ExtractJSON(packet []byte) (any, error) {
	root := (*document.DocumentSet)(doc).Root
	result, err := decoder.Decode(root, packet)
	if err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	jsonRes, err := json_out.BuildJSONWithSet((*document.DocumentSet)(doc), doc.Root, result)
	if err != nil {
		return nil, fmt.Errorf("failed to build JSON with set: %w", err)
	}

	return jsonRes, nil
}

// PacketType is a type alias for string, representing the name of the main packet type defined in the PDL document.
type PacketType string

// WithAs returns a new PacketType that represents the current packet type with an "as " prefix, which can be used in PDL
// source definitions to indicate that this packet type should be treated as an alias for another packet type during decoding.
func (pt PacketType) WithAs() PacketType {
	return "as " + pt
}

// String returns the string representation of the PacketType, which is simply the underlying string value.
func (pt PacketType) String() string {
	return string(pt)
}
