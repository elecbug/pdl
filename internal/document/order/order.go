package order

type ByteOrder string

const (
	BIG_ENDIAN    ByteOrder = "BIG_ENDIAN"
	LITTLE_ENDIAN ByteOrder = "LITTLE_ENDIAN"
)

type BitOrder string

const (
	MSB_FIRST BitOrder = "MSB_FIRST"
	LSB_FIRST BitOrder = "LSB_FIRST"
)
