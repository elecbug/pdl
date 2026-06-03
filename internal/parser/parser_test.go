package parser_test

import (
	"testing"

	"github.com/elecbug/pdl/internal/parser"
	"github.com/elecbug/pdl/internal/standard"
)

func TestParseForTCP(t *testing.T) {
	doc, err := parser.ParseString(standard.TCP_PDL)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if doc.PacketName != "TCP" {
		t.Fatalf("PacketName = %q, want TCP", doc.PacketName)
	}

	if len(doc.Defs) == 0 {
		t.Fatal("Defs is empty")
	}

	if len(doc.Outs) == 0 {
		t.Fatal("Outputs is empty")
	}
}
