package decoder

import (
	"testing"

	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/document/order"
)

func newTestContext(data []byte) *decodeContext {
	return &decodeContext{
		doc: &document.Document{
			ByteOrder: order.BIG_ENDIAN,
		},
		data:   data,
		values: map[string]Value{},
		vars:   map[string]int64{},
	}
}

func TestDecodeDef(t *testing.T) {
	t.Run("decodes with length", func(t *testing.T) {
		ctx := newTestContext([]byte{0x12, 0x34})

		def := document.Def{
			Name:      "f",
			From:      document.NumberExpr{Value: 0},
			UseLength: true,
			Length:    document.NumberExpr{Value: 8},
		}

		if err := ctx.decodeDef(def); err != nil {
			t.Fatalf("decodeDef failed: %v", err)
		}

		v, ok := ctx.values["f"]
		if !ok {
			t.Fatal("decoded value f not found")
		}

		if v.Len != 8 {
			t.Fatalf("Len = %d, want 8", v.Len)
		}

		if v.UInt != 0x12 {
			t.Fatalf("UInt = %d, want %d", v.UInt, 0x12)
		}
	})

	t.Run("decodes with to", func(t *testing.T) {
		ctx := newTestContext([]byte{0b10101010})

		def := document.Def{
			Name:  "f",
			From:  document.NumberExpr{Value: 2},
			UseTo: true,
			To:    document.NumberExpr{Value: 5},
		}

		if err := ctx.decodeDef(def); err != nil {
			t.Fatalf("decodeDef failed: %v", err)
		}

		if got := ctx.values["f"].Len; got != 4 {
			t.Fatalf("Len = %d, want 4", got)
		}
	})

	t.Run("decodes to end", func(t *testing.T) {
		ctx := newTestContext([]byte{0xff, 0x00})

		def := document.Def{
			Name:  "tail",
			From:  document.NumberExpr{Value: 8},
			UseTo: true,
			To:    document.EndExpr{},
		}

		if err := ctx.decodeDef(def); err != nil {
			t.Fatalf("decodeDef failed: %v", err)
		}

		if got := ctx.values["tail"].Len; got != 8 {
			t.Fatalf("Len = %d, want 8", got)
		}
	})

	t.Run("fails when range exceeds packet", func(t *testing.T) {
		ctx := newTestContext([]byte{0x00})

		def := document.Def{
			Name:      "overflow",
			From:      document.NumberExpr{Value: 6},
			UseLength: true,
			Length:    document.NumberExpr{Value: 4},
		}

		if err := ctx.decodeDef(def); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("fails when from is negative", func(t *testing.T) {
		ctx := newTestContext([]byte{0x00})

		def := document.Def{
			Name:      "neg",
			From:      document.NumberExpr{Value: -1},
			UseLength: true,
			Length:    document.NumberExpr{Value: 1},
		}

		if err := ctx.decodeDef(def); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("fails when length or to is missing", func(t *testing.T) {
		ctx := newTestContext([]byte{0x00})

		def := document.Def{
			Name: "bad",
			From: document.NumberExpr{Value: 0},
		}

		if err := ctx.decodeDef(def); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestEvalExpr(t *testing.T) {
	ctx := newTestContext([]byte{0x00, 0x00})
	ctx.vars["x"] = 7
	ctx.values["f"] = Value{UInt: 11}

	tests := []struct {
		name    string
		expr    document.Expr
		want    int64
		wantErr bool
	}{
		{
			name: "number",
			expr: document.NumberExpr{Value: 3},
			want: 3,
		},
		{
			name: "identifier",
			expr: document.IdentExpr{Name: "x"},
			want: 7,
		},
		{
			name: "field reference",
			expr: document.FieldValueExpr{Name: "f"},
			want: 11,
		},
		{
			name: "end expression",
			expr: document.EndExpr{},
			want: 15,
		},
		{
			name: "binary expression",
			expr: document.BinaryExpr{
				Op:   "+",
				Left: document.NumberExpr{Value: 2},
				Right: document.BinaryExpr{
					Op:    "*",
					Left:  document.NumberExpr{Value: 3},
					Right: document.NumberExpr{Value: 4},
				},
			},
			want: 14,
		},
		{
			name:    "undefined variable",
			expr:    document.IdentExpr{Name: "missing"},
			wantErr: true,
		},
		{
			name:    "undecoded field",
			expr:    document.FieldValueExpr{Name: "missing"},
			wantErr: true,
		},
		{
			name: "division by zero",
			expr: document.BinaryExpr{
				Op:    "/",
				Left:  document.NumberExpr{Value: 10},
				Right: document.NumberExpr{Value: 0},
			},
			wantErr: true,
		},
		{
			name: "unknown operator",
			expr: document.BinaryExpr{
				Op:    "%",
				Left:  document.NumberExpr{Value: 10},
				Right: document.NumberExpr{Value: 3},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ctx.evalExpr(tc.expr)
			if (err != nil) != tc.wantErr {
				t.Fatalf("evalExpr() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && got != tc.want {
				t.Fatalf("evalExpr() = %d, want %d", got, tc.want)
			}
		})
	}
}
