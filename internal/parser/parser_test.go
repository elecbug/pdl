package parser_test

import (
	"testing"

	"github.com/elecbug/pdl/internal/parser"
)

const tcpPDL = `
packet TCP

set mode BIG_ENDIAN MSB_FIRST

var {
    fixed_header_bits = 160
}

def {
    src_port       from 0   length 16
    dst_port       from 16  length 16
    seq            from 32  to     63
    ack            from 64  length 32
    data_offset    from 96  length 4
    reserved       from 100 length 3
    ns             from 103 length 1
    flags          from 104 length 8
    window         from 112 length 16
    checksum       from 128 length 16
    urgent_pointer from 144 length 16
    options        from fixed_header_bits length (*data_offset * 32 - fixed_header_bits)
    payload        from (*data_offset * 32) to end
}

out json {
    src_port source.port DEC
    dst_port destination.port DEC

    flags<6> flags.syn {
        0 : false
        1 : true
    }
}
`

func TestParseTCP(t *testing.T) {
	doc, err := parser.ParseString(tcpPDL)
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
