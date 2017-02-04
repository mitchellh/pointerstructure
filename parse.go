package pointerstructure

import (
	"fmt"
	"strings"
)

// Parse parses a pointer from the input string. The input string
// is expected to follow the format specified by RFC 6901: '/'-separated
// parts. Each part can contain escape codes to contain '/' or '~'.
func Parse(input string) (*Pointer, error) {
	// Special case the empty case
	if input == "" {
		return &Pointer{}, nil
	}

	// We expect the first character to be "/"
	if input[0] != '/' {
		return nil, fmt.Errorf(
			"parse Go pointer %q: first char must be '/'", input)
	}

	// Trim out the first slash so we don't have to +1 every index
	input = input[1:]

	// Parse out all the parts
	var parts []string
	lastSlash := -1
	for i, r := range input {
		if r == '/' {
			parts = append(parts, input[lastSlash+1:i])
			lastSlash = i
		}
	}

	// Add last part
	parts = append(parts, input[lastSlash+1:])

	// Process each part for string replacement
	for i, p := range parts {
		// Replace ~1 followed by ~0 as specified by the RFC
		parts[i] = strings.Replace(
			strings.Replace(p, "~1", "/", -1), "~0", "~", -1)
	}

	return &Pointer{Parts: parts}, nil
}
