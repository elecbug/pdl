package pdl

// Source represents a raw PDL source string that can be used to create a PDL document for decoding packets and extracting JSON output.
// It provides a String method to retrieve the underlying string value of the source code.
type Source struct {
	src string
}

// NewSource creates a new Source instance from a given string, which represents the raw PDL source code that can be used
// to create a PDL document for decoding packets and extracting JSON output.
func NewSource(src string) Source {
	return Source{
		src: src,
	}
}

// String returns the string representation of the Source, which is simply the underlying string value of the source code.
func (s Source) String() string {
	return s.src
}
