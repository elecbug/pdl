package pdl

import (
	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/extractor/json_out"
	"github.com/elecbug/pdl/internal/parser"
)

type Document document.Document

func GenerateDocument(src string) (*Document, error) {
	doc, err := parser.ParseString(src)
	if err != nil {
		return nil, err
	}
	return (*Document)(doc), nil
}

func SerializeDocument(d *Document) ([]byte, error) {
	return (*document.Document)(d).Serialize()
}

func DeserializeDocument(data []byte) (*Document, error) {
	doc, err := document.Deserialize(data)
	if err != nil {
		return nil, err
	}

	return (*Document)(doc), nil
}

func ExtractJSON(d *Document, packet []byte) (any, error) {
	result, err := decoder.Decode((*document.Document)(d), packet)
	if err != nil {
		return nil, err
	}

	jsonResult, err := json_out.BuildJSON((*document.Document)(d), result)
	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}
