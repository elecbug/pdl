package extractor_test

import (
	"testing"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/extractor/json_out"
	"github.com/elecbug/pdl/internal/parser"
)

const tcpPDL = `
packet TCP

set mode BIG_ENDIAN MSB_FIRST

var {
    fixed_header_bits = 160
}

def {
    src_port       from 0     length 16
    dst_port       from 16    length 16
    seq            from 32    length 32
    ack            from 64    length 32
    data_offset    from 96    length 4
    reserved       from 100   length 3
    ns             from 103   length 1
    flags          from 104   length 8
	test_flags	   from 104+4 length 4
    window         from 112   length 16
    checksum       from 128   length 16
    urgent_pointer from 144   length 16
    options        from fixed_header_bits length (*data_offset * 32 - fixed_header_bits)
    payload        from (*data_offset * 32) to end
}

out json {
    src_port       source.port           DEC
    dst_port       destination.port      DEC
    seq            sequence_number       DEC
    ack            acknowledgment_number DEC
    data_offset    header.length_words   DEC
    reserved       header.reserved       BIN
    ns             header.ns             BOOL
    window         window_size           DEC
    checksum       checksum              HEX
    urgent_pointer urgent_pointer        DEC
    options        options               HEX
    payload        payload               HEX

	flags<4> flags.psh {
		0 : false
		1 : true
	}

	flags<5> flags.rst {
		0 : false
		1 : true
	}

    flags<6> flags.syn {
        0 : false
        1 : true
    }

    flags<7> flags.fin {
        0 : false
        1 : true
    }

	test_flags test_flags {
		0b0000 : "none"
		0b0001 : "fin"
		0b0010 : "syn"
		0b0011 : "fin+syn"
		0b0100 : "rst"
		0b0101 : "fin+rst"
		0b0110 : "syn+rst"
		0b0111 : "fin+syn+rst"
		0b1000 : "psh"
		0b1001 : "fin+psh"
		0b1010 : "syn+psh"
		0b1011 : "fin+syn+psh"
		0b1100 : "rst+psh"
		0b1101 : "fin+rst+psh"
		0b1110 : "syn+rst+psh"
		0b1111 : "all"
	}
}
`

func TestBuildJSONTCP(t *testing.T) {
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
		0x0a,
		0x20, 0x00,
		0xab, 0xcd,
		0x00, 0x00,
		0xde, 0xad, 0xbe, 0xef,
	}

	result, err := decoder.Decode(doc, packet)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	obj, err := json_out.BuildJSON(doc, result)
	if err != nil {
		t.Fatalf("BuildJSON failed: %v", err)
	}

	root, ok := obj.(map[string]any)
	if !ok {
		t.Fatalf("json root type = %T, want map[string]any", obj)
	}

	source := root["source"].(map[string]any)
	if source["port"] != uint64(80) {
		t.Fatalf("source.port = %v, want 80", source["port"])
	}

	destination := root["destination"].(map[string]any)
	if destination["port"] != uint64(443) {
		t.Fatalf("destination.port = %v, want 443", destination["port"])
	}

	flags := root["flags"].(map[string]any)
	if flags["syn"] != true {
		t.Fatalf("flags.syn = %v, want true", flags["syn"])
	}

	if flags["fin"] != false {
		t.Fatalf("flags.fin = %v, want false", flags["fin"])
	}

	if flags["psh"] != true {
		t.Fatalf("flags.psh = %v, want true", flags["psh"])
	}

	if flags["rst"] != false {
		t.Fatalf("flags.rst = %v, want false", flags["rst"])
	}

	if root["test_flags"] != "syn+psh" {
		t.Fatalf("test_flags = %v, want syn+psh", root["test_flags"])
	}

	if root["payload"] != "DEADBEEF" {
		t.Fatalf("payload = %v, want DEADBEEF", root["payload"])
	}
}
