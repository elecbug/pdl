package decoder

import (
	"fmt"

	"github.com/elecbug/pdl/internal/document"
)

// Value represents a decoded field value, including its name, raw bits, length in bits, unsigned
// integer value, and mode (e.g., "hex", "dec").
type Value struct {
	// The name of the decoded field, as defined in the document.
	Name string
	// The raw bits of the decoded field, extracted from the input data.
	Bits []byte
	// The length of the decoded field in bits.
	Len int64
	// The unsigned integer value of the decoded field, derived from the raw bits.
	UInt uint64
	// The output format mode for the decoded field, such as "hex", "dec",
	// "bin", or "bool".
	Mode string
}

// Result represents the outcome of the decoding process, containing a map of decoded field values
// keyed by their names.
type Result struct {
	// A map of decoded field values, keyed by field name.
	Values map[string]Value
}

// Decode takes a document and input data, and returns the decoded result or an error if the decoding
// process fails. It initializes a DecodeContext, evaluates any variables defined in the document,
// and then decodes each field definition according to the specified rules.
func Decode(doc *document.Document, data []byte) (*Result, error) {
	ctx := &decodeContext{
		doc:    doc,
		data:   data,
		values: make(map[string]Value),
		vars:   make(map[string]int64),
	}

	for _, v := range doc.Vars {
		value, err := ctx.evalExpr(v.Expr)
		if err != nil {
			return nil, fmt.Errorf("var %s: %w", v.Name, err)
		}

		ctx.vars[v.Name] = value
	}

	for _, def := range doc.Defs {
		if err := ctx.decodeDef(def); err != nil {
			return nil, err
		}
	}

	return &Result{
		Values: ctx.values,
	}, nil
}
