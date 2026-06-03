package formatter_test

import (
	"testing"

	"github.com/elecbug/pdl/internal/decoder"
	"github.com/elecbug/pdl/internal/document/order"
	"github.com/elecbug/pdl/internal/formatter"
)

func TestFormatValue(t *testing.T) {
	v := decoder.Value{
		Name: "test",
		Bits: []byte{0x12, 0x34},
		Len:  16,
		UInt: 0x1234,
	}

	tests := []struct {
		format   string
		expected any
	}{
		{"DEC", uint64(0x1234)},
		{"HEX", "1234"},
		{"BIN", "0001001000110100"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			out, err := formatter.FormatValue(v, tt.format)
			if err != nil {
				t.Fatalf("FormatValue failed: %v", err)
			}
			if out != tt.expected {
				t.Fatalf("FormatValue = %v, want %v", out, tt.expected)
			}
		})
	}
}

func TestGetBit(t *testing.T) {
	v := decoder.Value{
		Name: "test",
		Bits: []byte{0b10101010},
		Len:  8,
		UInt: 0b10101010,
	}

	tests := []struct {
		name     string
		bitOrder order.BitOrder
		expected []uint64
	}{
		{"MSB_FIRST", order.MSB_FIRST, []uint64{1, 0, 1, 0, 1, 0, 1, 0}},
		{"LSB_FIRST", order.LSB_FIRST, []uint64{0, 1, 0, 1, 0, 1, 0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < int(v.Len); i++ {
				bit, err := formatter.GetBit(v, i, tt.bitOrder)
				if err != nil {
					t.Fatalf("GetBit failed: %v", err)
				}
				if bit != tt.expected[i] {
					t.Fatalf("GetBit(%d) = %d, want %d", i, bit, tt.expected[i])
				}
			}
		})
	}
}

func TestConvertMappedValue(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"true", true},
		{"false", false},
		{"other", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			out := formatter.ConvertMappedValue(tt.input)
			if out != tt.expected {
				t.Fatalf("ConvertMappedValue(%q) = %v, want %v", tt.input, out, tt.expected)
			}
		})
	}
}
