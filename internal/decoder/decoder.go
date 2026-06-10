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

// ArrayValue represents a decoded array of fields, including the name of the array, the packet type of each item,
type ArrayValue struct {
	// The name of the array field, as defined in the document.
	Name string
	// The packet type of each item in the array, which should correspond to a packet defined in the document.
	Packet string
	// The list of decoded items in the array, where each item includes its raw bits, length in bits, and decoded result.
	Items []ArrayItem
}

// ArrayItem represents a single item in a decoded array, including its raw bits, length in bits, and decoded result.
type ArrayItem struct {
	// The raw bits of the decoded item, extracted from the input data.
	Bits []byte
	// The length of the decoded item in bits.
	Len int64
	// The decoded result of the item, which is a mapping of field names to their corresponding decoded values.
	Result *Result
}

// Result represents the overall result of the decoding process, including a mapping of field names to their decoded values,
// a mapping of array field names to their decoded array values, and the total number of bits consumed during decoding.
type Result struct {
	// A mapping of field names to their corresponding decoded values, where each value includes the raw bits, length in bits, unsigned integer value, and output mode.
	Values map[string]Value
	// A mapping of array field names to their corresponding decoded array values, where each array value includes the packet type and list of decoded items.
	Arrays map[string]ArrayValue
	// The total number of bits consumed during the decoding process, which can be used for tracking the position in the input data.
	ConsumedBits int64
}

// Decode takes a Document and input data as a byte slice, and returns the decoded Result or an error if decoding fails.
func Decode(doc *document.Document, data []byte) (*Result, error) {
	return decode(nil, doc, data)
}

// DecodeWithSet takes a DocumentSet, a Document, and input data as a byte slice, and returns the decoded Result or an error if decoding fails.
func DecodeWithSet(set *document.DocumentSet, doc *document.Document, data []byte) (*Result, error) {
	return decode(set, doc, data)
}

// decode is the internal function that performs the actual decoding logic, using the provided DocumentSet, Document, and input data.
func decode(set *document.DocumentSet, doc *document.Document, data []byte) (*Result, error) {
	ctx := &decodeContext{
		set:    set,
		doc:    doc,
		data:   data,
		values: make(map[string]Value),
		arrays: make(map[string]ArrayValue),
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
		Values:       ctx.values,
		Arrays:       ctx.arrays,
		ConsumedBits: ctx.consumedBits,
	}, nil
}
