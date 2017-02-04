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

// Set sets the value at the given pointer.
//
// This is a shorthand for calling Parse on the pointer and then calling Set
// on that result. An error will be returned if the value cannot be found or
// there is an error with the format of pointer.
//
// Set returns the complete document, which might change if the pointer value
// points to the root ("").
func Set(doc interface{}, pointer string, value interface{}) (interface{}, error) {
	p, err := Parse(pointer)
	if err != nil {
		return nil, err
	}

	return p.Set(doc, value)
}

// Parent returns a pointer to the parent element of this pointer.
//
// If Pointer represents the root (empty parts), a pointer representing
// the root is returned. Therefore, to check for the root, IsRoot() should be
// called.
func (p *Pointer) Parent() *Pointer {
	// If this is root, then we just return a new root pointer. We allocate
	// a new one though so this can still be modified.
	if p.IsRoot() {
		return &Pointer{}
	}

	parts := make([]string, len(p.Parts)-1)
	copy(parts, p.Parts[:len(p.Parts)-1])
	return &Pointer{
		Parts: parts,
	}
}

// IsRoot returns true if this pointer represents the root document.
func (p *Pointer) IsRoot() bool {
	return len(p.Parts) == 0
}
