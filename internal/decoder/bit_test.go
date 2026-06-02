package decoder

import (
	"slices"
	"testing"

	"github.com/elecbug/pdl/internal/document/order"
)

func TestExtractBits(t *testing.T) {
	mainData := []byte{
		0b00010000, 0b00100000, 0b01000000, 0b10000000,
		0b00000001, 0b00000010, 0b00000100, 0b00001000,
		0b00010001, 0b00100010, 0b01000100, 0b10001000,
		0b00010010, 0b00100100, 0b01001000, 0b10010000,
	}

	testCases := []struct {
		name     string
		data     []byte
		from     int64
		length   int64
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Extract bits 0-15",
			data:     mainData,
			from:     0,
			length:   16,
			expected: []byte{0b00010000, 0b00100000},
			wantErr:  false,
		},
		{
			name:     "Extract bits 4-15",
			data:     mainData,
			from:     4,
			length:   12,
			expected: []byte{0b00000010, 0b00000000},
			wantErr:  false,
		},
		{
			name:     "Extract bits 8-27",
			data:     mainData,
			from:     8,
			length:   20,
			expected: []byte{0b00100000, 0b01000000, 0b10000000},
			wantErr:  false,
		},
		{
			name:     "Extract bits 60-67",
			data:     mainData,
			from:     60,
			length:   8,
			expected: []byte{0b10000001},
			wantErr:  false,
		},
		{
			name:     "Extract bits 96-106",
			data:     mainData,
			from:     96,
			length:   11,
			expected: []byte{0b00010010, 0b00100000},
			wantErr:  false,
		},
		{
			name:     "Extract bits 0-127",
			data:     mainData,
			from:     0,
			length:   128,
			expected: mainData,
			wantErr:  false,
		},
		{
			name:    "Invalid range: negative from",
			data:    mainData,
			from:    -1,
			length:  8,
			wantErr: true,
		},
		{
			name:    "Invalid range: negative length",
			data:    mainData,
			from:    0,
			length:  -1,
			wantErr: true,
		},
		{
			name:    "Invalid range: exceeds data size",
			data:    mainData,
			from:    120,
			length:  16,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			extracted, err := extractBits(tc.data, tc.from, tc.length)
			if (err != nil) != tc.wantErr {
				t.Fatalf("extractBits() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil && !tc.wantErr && !slices.Equal(extracted, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, extracted)
			}
		})
	}
}

func TestBitsToUint(t *testing.T) {
	testCases := []struct {
		name      string
		bits      []byte
		bitLen    int64
		byteOrder order.ByteOrder
		want      uint64
		wantErr   bool
	}{
		{
			name:      "Test big endian 16 bits",
			bits:      []byte{0b00010000, 0b00100000},
			bitLen:    16,
			byteOrder: order.BIG_ENDIAN,
			want:      0b0001000000100000,
		},
		{
			name:      "Test little endian 16 bits",
			bits:      []byte{0b00010000, 0b00100000},
			bitLen:    16,
			byteOrder: order.LITTLE_ENDIAN,
			want:      0b0010000000010000,
		},
		{
			name:      "Test big	endian 20 bits",
			bits:      []byte{0b00010000, 0b00100000, 0b01000000},
			bitLen:    20,
			byteOrder: order.BIG_ENDIAN,
			want:      0b00010000001000000100,
		},
		{
			name:      "Test little endian 20 bits",
			bits:      []byte{0b00010000, 0b00100000, 0b01000000},
			bitLen:    20,
			byteOrder: order.LITTLE_ENDIAN,
			wantErr:   true, // Not byte-aligned
		},
		{
			name:      "Test big endian 64 bits",
			bits:      []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
			bitLen:    64,
			byteOrder: order.BIG_ENDIAN,
			want:      0x0123456789abcdef,
		},
		{
			name:      "Test little endian 64 bits",
			bits:      []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
			bitLen:    64,
			byteOrder: order.LITTLE_ENDIAN,
			want:      0xefcdab8967452301,
		},
		{
			name:      "Test too long bit length",
			bits:      []byte{0x00},
			bitLen:    65,
			byteOrder: order.BIG_ENDIAN,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bitsToUint(tc.bits, tc.bitLen, tc.byteOrder)
			if (err != nil) != tc.wantErr {
				t.Fatalf("bitsToUint() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil && got != tc.want {
				t.Errorf("bitsToUint() = %v, want %v", got, tc.want)
			}
		})
	}
}
