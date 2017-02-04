// Package pointerstructure provides functions for identifying a specific
// value within any Go structure using a string syntax.
//
// The syntax used is based on JSON Pointer (RFC 6901).
package pointerstructure

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// Pointer represents a pointer to a specific value. You can construct
// a pointer manually or use Parse.
type Pointer struct {
	// Parts are the pointer parts. No escape codes are processed here.
	// The values are expected to be exact. If you have escape codes, use
	// the Parse functions.
	Parts []string
}

// Get reads the value out of the total value v.
func (p *Pointer) Get(v interface{}) (interface{}, error) {
	// fast-path the empty address case to avoid reflect.ValueOf below
	if len(p.Parts) == 0 {
		return v, nil
	}

	var err error
	currentVal := reflect.ValueOf(v)
	for i, part := range p.Parts {
		for currentVal.Kind() == reflect.Interface {
			currentVal = currentVal.Elem()
		}

		switch currentVal.Kind() {
		case reflect.Map:
			currentVal, err = p.getMap(part, currentVal)
			if err != nil {
				return nil, fmt.Errorf("%s at part %d: %s", p, i, err)
			}

		default:
			return nil, fmt.Errorf(
				"%s: at part %d, invalid value kind: %s", p, i, currentVal.Kind())
		}
	}

	return currentVal.Interface(), nil
}

func (p *Pointer) getMap(part string, m reflect.Value) (reflect.Value, error) {
	var zeroValue reflect.Value

	// Determine the key type so we can convert into it
	keyType := m.Type().Key()
	key := reflect.New(keyType)
	if err := mapstructure.WeakDecode(part, key.Interface()); err != nil {
		return zeroValue, fmt.Errorf(
			"couldn't convert key %q to type %s", part, keyType.String())
	}

	// Need to dereference the pointer since reflect.New always
	// creates a pointer.
	key = reflect.Indirect(key)

	// Verify that the key exists
	found := false
	for _, k := range m.MapKeys() {
		if k.Interface() == key.Interface() {
			found = true
			break
		}
	}
	if !found {
		return zeroValue, fmt.Errorf("couldn't find key %#v", key.Interface())
	}

	// Get the key
	return m.MapIndex(key), nil
}
