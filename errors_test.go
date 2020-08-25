package pointerstructure

import (
	"errors"
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
	if !errors.Is(err, ErrParse) {
		t.Fatalf("expected ErrParse in the error chain, but it was not")
	}
}

func TestErrNotFound(t *testing.T) {
	_, err := Get(structure, notFound)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound in the error chain, but it was not")
	}
}

func TestErrOutOfRange(t *testing.T) {
	_, err := Get(structure, outOfRange)
	if !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("expected ErrOutOfRange in the error chain, but it was not")
	}
}

func TestErrConvert(t *testing.T) {
	_, err := Set(structure, cantConvert, "test")
	if !errors.Is(err, ErrConvert) {
		t.Fatalf("expected ErrConvert in the error chain, but it was not")
	}
}

func TestErrInvalidKind(t *testing.T) {
	_, err := Get(structure, invalidKind)
	if !errors.Is(err, ErrInvalidKind) {
		t.Fatalf("expected ErrInvalidKind in the error chain, but it was not")
	}
}
