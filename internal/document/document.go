package document

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	"github.com/elecbug/pdl/internal/document/order"
)

type Document struct {
	PacketName string

	ByteOrder order.ByteOrder
	BitOrder  order.BitOrder

	Vars []Var
	Defs []Def
	Outs []Out
}

type Var struct {
	Name string
	Expr Expr
}

type Def struct {
	Name string

	From Expr

	Length Expr
	To     Expr

	UseLength bool
	UseTo     bool
}

type Out struct {
	Field  string
	Path   string
	Format string

	HasBitIndex bool
	BitIndex    int

	Map map[string]string
}

var inited = false

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
