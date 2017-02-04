// Package pointerstructure provides functions for identifying a specific
// value within any Go structure using a string syntax.
//
// The syntax used is based on JSON Pointer (RFC 6901).
package pointerstructure

// Pointer represents a pointer to a specific value. You can construct
// a pointer manually or use Parse.
type Pointer struct {
	// Parts are the pointer parts. No escape codes are processed here.
	// The values are expected to be exact. If you have escape codes, use
	// the Parse functions.
	Parts []string
}

// Get reads the value at the given pointer.
//
// This is a shorthand for calling Parse on the pointer and then calling Get
// on that result. An error will be returned if the value cannot be found or
// there is an error with the format of pointer.
func Get(value interface{}, pointer string) (interface{}, error) {
	p, err := Parse(pointer)
	if err != nil {
		return nil, err
	}

	return p.Get(value)
}
