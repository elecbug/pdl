package decoder_test

import (
	"testing"

	"github.com/elecbug/pdl/internal/decoder"
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
    seq            from 32  length 32
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
}
`

func TestDecodeTCP(t *testing.T) {
	doc, err := parser.ParseString(tcpPDL)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	packet := []byte{
		0x00, 0x50,
		0x01, 0xbb,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x50,
		0x02,
		0x20, 0x00,
		0xab, 0xcd,
		0x00, 0x00,
		0xde, 0xad, 0xbe, 0xef,
	}

	result, err := decoder.Decode(doc, packet)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	tests := map[string]uint64{
		"src_port":       80,
		"dst_port":       443,
		"seq":            1,
		"ack":            0,
		"data_offset":    5,
		"reserved":       0,
		"ns":             0,
		"flags":          2,
		"window":         8192,
		"checksum":       0xabcd,
		"urgent_pointer": 0,
	}

	for name, want := range tests {
		got, ok := result.Values[name]
		if !ok {
			t.Fatalf("missing decoded field %q", name)
		}

		if got.UInt != want {
			t.Fatalf("%s = %d, want %d", name, got.UInt, want)
		}
	}

	payload, ok := result.Values["payload"]
	if !ok {
		t.Fatal("missing payload")
	}

	if payload.Len != 32 {
		t.Fatalf("payload.Len = %d, want 32", payload.Len)
	}
}