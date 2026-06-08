package document

import (
	"github.com/elecbug/pdl/internal/document/order"
)

// Document represents the structured representation of a PDL document, containing the packet name,
// byte order, bit order, variables, field definitions, and output specifications.
type Document struct {
	// The name of the packet being defined in the document.
	PacketName string

	// The byte order (endianness), either BIG_ENDIAN or LITTLE_ENDIAN.
	ByteOrder order.ByteOrder
	// The bit order, either MSB_FIRST or LSB_FIRST.
	BitOrder order.BitOrder

	// A list of variables defined in the document, where each variable has a name and an associated
	// expression that can be evaluated during decoding.
	Vars []Var
	// A list of field definitions in the document, where each definition specifies how to extract a
	// field from the input data based on its position and length.
	Defs []Def
	// A list of output specifications in the document, where each output defines how to format and
	// present a decoded field value, including optional mapping for value formatting.
	Outs []Out
}

// DocumentSet represents a collection of documents, allowing for multiple packet definitions and a
// designated root document for decoding.
type DocumentSet struct {
	// A map of document names to their corresponding Document structures, enabling the organization
	Documents map[string]*Document
	// A reference to the root document that should be used for decoding input data,
	// which must be one of the documents in the Documents map.
	Root *Document
}

// Var represents a variable defined in the document, consisting of a name and an associated expression
type Var struct {
	// The name of the variable, which can be used in expressions within the document.
	Name string
	// The expression associated with this variable.
	Expr Expr
}

// Def represents a field definition in the document, specifying how to extract a field from the input
type Def struct {
	// The field name used as the key in decoded results.
	Name string

	// The expression specifying the starting bit position of the field within the input data.
	From Expr
	// The expression specifying the length of the field in bits, if UseLength is true.
	Length Expr
	// The ending bit position expression when UseTo is true.
	To Expr

	// UseLength indicates whether the Length expression should be used to determine the field's length.
	UseLength bool
	// UseTo indicates whether the To expression should be used to determine the field's ending position.
	UseTo bool
}

// Out represents an output specification in the document, defining how to format and present a decoded
type Out struct {
	// The name of the field to output, which should correspond to a field defined in the document.
	Field string
	// The destination JSON path.
	Path string
	// The output format, such as "hex", "dec", "bin", or "bool".
	Format string

	// Whether this output maps a single bit selected by BitIndex.
	HasBitIndex bool
	// The bit index used when HasBitIndex is true.
	BitIndex int

	// Map is an optional mapping of decoded values to their corresponding string representations
	// which can be used for formatting the output.
	Map map[string]string
	// MapDefault is an optional default value to use when a decoded value does not have a corresponding entry in the Map.
	MapDefault *string

	// AsPacket is an optional field that, if set, indicates that the output should be treated as a nested packet
	AsPacket string
}
