package pointerstructure

import (
	"fmt"
	"reflect"
)

// Set writes a value v to the pointer p in structure s.
//
// The structures s must have non-zero values set up to this pointer.
// For example, if setting "/bob/0/name", then "/bob/0" must be set already.
//
// The returned value is potentially a new value if this pointer represents
// the root document. Otherwise, the returned value will always be s.
func (p *Pointer) Set(s, v interface{}) (interface{}, error) {
	// if we represent the root doc, return that
	if len(p.Parts) == 0 {
		return v, nil
	}

	// Save the original since this is going to be our return value
	originalS := s

	// Get the parent value
	var err error
	s, err = p.Parent().Get(s)
	if err != nil {
		return nil, err
	}

	// Map for lookup of getter to call for type
	funcMap := map[reflect.Kind]func(string, reflect.Value, reflect.Value) error{
		reflect.Array: p.setSlice,
		reflect.Map:   p.setMap,
		reflect.Slice: p.setSlice,
	}

	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	for val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	f, ok := funcMap[val.Kind()]
	if !ok {
		return nil, fmt.Errorf("set %s: invalid value kind: %s", p, val.Kind())
	}
	if err := f(p.Parts[len(p.Parts)-1], val, reflect.ValueOf(v)); err != nil {
		return nil, fmt.Errorf("set %s: %s", p, err)
	}

	return originalS, nil
}

func (p *Pointer) setMap(part string, m, value reflect.Value) error {
	key, err := coerce(reflect.ValueOf(part), m.Type().Key())
	if err != nil {
		return err
	}

	elem, err := coerce(value, m.Type().Elem())
	if err != nil {
		return err
	}

	// Set the key
	m.SetMapIndex(key, elem)
	return nil
}

func (p *Pointer) setSlice(part string, s, value reflect.Value) error {
	// Coerce the key to an int
	idxVal, err := coerce(reflect.ValueOf(part), reflect.TypeOf(42))
	if err != nil {
		return err
	}
	idx := int(idxVal.Int())

	// Verify we're within bounds
	if idx < 0 || idx >= s.Len() {
		return fmt.Errorf(
			"index %d is out of range (length = %d)", idx, s.Len())
	}

	// Coerce the value
	value, err = coerce(value, s.Type().Elem())
	if err != nil {
		return err
	}

	// Set the key
	s.Index(idx).Set(value)
	return nil
}
