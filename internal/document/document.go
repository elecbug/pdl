package document

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	"github.com/elecbug/pdl/internal/document/order"
)

// Document represents the structured representation of a PDL document, containing the packet name,
// byte order, bit order, variables, field definitions, and output specifications.
type Document struct {
	// The name of the packet being defined in the document.
	PacketName string

	// The byte order (endianness) used in the document, which can be either BIG_ENDIAN or LITTLE_ENDIAN.
	ByteOrder order.ByteOrder
	// The bit order used in the document, which can be either MSB_FIRST or LSB_FIRST.
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

// Var represents a variable defined in the document, consisting of a name and an associated expression
type Var struct {
	// The name of the variable, which can be used in expressions within the document.
	Name string
	// The expression associated with the variable, which can be evaluated during decoding to produce a value.
	Expr Expr
}

// Def represents a field definition in the document, specifying how to extract a field from the input
type Def struct {
	// The name of the field being defined, which will be used as the key for the decoded value in the result.
	Name string

	// The expression specifying the starting bit position of the field within the input data.
	From Expr
	// The expression specifying the length of the field in bits, if UseLength is true.
	Length Expr
	// The expression specifying the ending bit position of the field within the input data, if UseTo is true.
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
	// The format to use when presenting the decoded field value, such as "hex", "dec", "bin", or "bool".
	Path string
	// The format string to use when formatting the decoded field value, which can include placeholders for the value.
	Format string

	// Map is an optional mapping of decoded values to their corresponding string representations,
	// which can be used for formatting the output.
	HasBitIndex bool
	// BitIndex specifies the index of the bit to check when HasBitIndex is true, allowing for conditional
	// formatting based on the value of a specific bit within the decoded field.
	BitIndex int

	// Map is an optional mapping of decoded values to their corresponding string representations
	// which can be used for formatting the output.
	Map map[string]string
}

// inited indicates whether the gob package has been initialized with the necessary type registrations
// for encoding and decoding Document instances.
var inited = false

// initGob initializes the gob package by registering the necessary types for encoding and decoding
// Document instances. This function is called before any serialization or deserialization operations
// to ensure that the gob encoder and decoder can properly handle the custom types used in the Document struct.
func initGob() {
	if !inited {
		gob.Register(NumberExpr{})
		gob.Register(IdentExpr{})
		gob.Register(FieldValueExpr{})
		gob.Register(EndExpr{})
		gob.Register(BinaryExpr{})
		inited = true
	}
}

// Serialize encodes the Document instance into a byte slice using gob encoding and base64 encoding.
// It first initializes the gob package, then encodes the Document instance into a buffer, and finally
// converts the encoded bytes to a base64 string before returning it as a byte slice.
func (d *Document) Serialize() ([]byte, error) {
	initGob()

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(d); err != nil {
		return nil, err
	}

	base64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return []byte(base64), nil
}

// Deserialize takes a byte slice containing a base64-encoded gob representation of a Document instance,
// decodes it, and returns the resulting Document instance or an error if the deserialization process fails.
// It first initializes the gob package, then decodes the base64 string into a buffer, and finally
// decodes the gob data from the buffer into a Document instance.
func Deserialize(data []byte) (*Document, error) {
	initGob()

	var doc Document

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	dec := gob.NewDecoder(bytes.NewReader(decoded))
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}

	return &doc, nil
}
