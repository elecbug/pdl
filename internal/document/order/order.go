package order

// ByteOrder represents the byte order (endianness) used in the document.
type ByteOrder string

const (
	// BIG_ENDIAN indicates that the most significant byte is stored at the smallest memory address.
	BIG_ENDIAN ByteOrder = "BIG_ENDIAN"
	// LITTLE_ENDIAN indicates that the least significant byte is stored at the smallest memory address.
	LITTLE_ENDIAN ByteOrder = "LITTLE_ENDIAN"
)

// BitOrder represents the bit order used in the document.
type BitOrder string

const (
	// MSB_FIRST indicates that the most significant bit is processed first within a byte.
	MSB_FIRST BitOrder = "MSB_FIRST"
	// LSB_FIRST indicates that the least significant bit is processed first within a byte.
	LSB_FIRST BitOrder = "LSB_FIRST"
)
