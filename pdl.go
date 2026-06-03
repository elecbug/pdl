package pdl

import (
	"fmt"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/extractor/json_out"
	"github.com/elecbug/pdl/internal/parser"
)

type Document document.DocumentSet

func GenerateDocument(root string, srcs ...string) (*Document, error) {
	doc, err := parser.ParseWithMultiSources(root, srcs...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	return (*Document)(doc), nil
}

func ExtractJSON(doc *Document, packet []byte) (any, error) {
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
