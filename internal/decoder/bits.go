package decoder

import (
	"fmt"

	"github.com/elecbug/pdl/internal/ast"
)

func extractBits(data []byte, from int64, length int64) []byte {
	out := make([]byte, (length+7)/8)

	for i := int64(0); i < length; i++ {
		srcBit := from + i

		srcByteIdx := srcBit / 8
		srcBitIdx := srcBit % 8

		// Network-style bit order inside byte:
		// bit 0 is MSB.
		bit := (data[srcByteIdx] >> (7 - srcBitIdx)) & 1

		dstByteIdx := i / 8
		dstBitIdx := i % 8

		out[dstByteIdx] |= bit << (7 - dstBitIdx)
	}

	return out
}

func bitsToUint(bits []byte, bitLen int64, byteOrder ast.ByteOrder) (uint64, error) {
	if bitLen > 64 {
		return 0, fmt.Errorf("cannot convert field longer than 64 bits to uint: %d", bitLen)
	}

	var u uint64

	switch byteOrder {
	case ast.BIG_ENDIAN:
		for i := int64(0); i < bitLen; i++ {
			byteIdx := i / 8
			bitIdx := i % 8

			bit := (bits[byteIdx] >> (7 - bitIdx)) & 1
			u = (u << 1) | uint64(bit)
		}

	case ast.LITTLE_ENDIAN:
		// MVP: reverse by byte if byte-aligned.
		// Bit-level little endian can be refined later.
		if bitLen%8 != 0 {
			return 0, fmt.Errorf("little endian for non-byte-aligned field is not supported yet")
		}

		byteLen := bitLen / 8
		for i := int64(0); i < byteLen; i++ {
			u |= uint64(bits[i]) << (8 * i)
		}

	default:
		return 0, fmt.Errorf("unknown byte order %q", byteOrder)
	}

	return u, nil
}
