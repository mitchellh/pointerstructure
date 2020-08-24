// +build !go1.13

package pointerstructure

import (
	"strings"
	"testing"
)

var structure interface{} = map[string]interface{}{
	"foo": map[string]interface{}{
		"bar": map[string]interface{}{
			"baz": []int{3, 0, 1, 5},
		},
		"quxx": "test",
	},
}

const (
	unparsable  = "foo/bar"
	notFound    = "/foo/baz/bar"
	outOfRange  = "/foo/bar/baz/5"
	cantConvert = "/foo/bar/baz/-"
	invalidKind = "/foo/quxx/test"
)

func TestErrParse(t *testing.T) {
	_, err := Parse(unparsable)
	if !strings.Contains(err.Error(), ErrParse.Error()) {
		t.Fatalf("expected ErrParse in the error chain, but it was not")
	}
}

func TestErrNotFound(t *testing.T) {
	_, err := Get(structure, notFound)
	if !strings.Contains(err.Error(), ErrNotFound.Error()) {
		t.Fatalf("expected ErrNotFound in the error chain, but it was not")
	}
}

func TestErrOutOfRange(t *testing.T) {
	_, err := Get(structure, outOfRange)
	if !strings.Contains(err.Error(), ErrOutOfRange.Error()) {
		t.Fatalf("expected ErrOutOfRange in the error chain, but it was not")
	}
}

func TestErrConvert(t *testing.T) {
	_, err := Set(structure, cantConvert, "test")
	if !strings.Contains(err.Error(), ErrConvert.Error()) {
		t.Fatalf("expected ErrConvert in the error chain, but it was not")
	}
}

func TestErrInvalidKind(t *testing.T) {
	_, err := Get(structure, invalidKind)
	if !strings.Contains(err.Error(), ErrInvalidKind.Error()) {
		t.Fatalf("expected ErrInvalidKind in the error chain, but it was not")
	}
}
