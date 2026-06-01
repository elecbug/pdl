package decoder

import (
	"fmt"

	"github.com/elecbug/pdl/internal/document"
)

type Value struct {
	Name string

	Bits []byte
	Len  int64

	UInt uint64

	Mode string
}

type Result struct {
	Values map[string]Value
}

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
