package ast

import (
	"bytes"
	"encoding/gob"
)

type Document struct {
	PacketName string

	ByteOrder ByteOrder
	BitOrder  BitOrder

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

	BitIndex *int
	Map      map[string]string
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

	return buf.Bytes(), nil
}

func Deserialize(data []byte) (*Document, error) {
	initGob()

	var doc Document

	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}

	return &doc, nil
}
