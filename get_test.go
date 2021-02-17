package pointerstructure

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPointerGet_Hook(t *testing.T) {
	type embedded struct {
		S interface{}
		V *embedded
	}
	type embedded2 struct {
		S2 interface{}
	}

	hookForInterface := func(v reflect.Value) reflect.Value {
		if !v.CanInterface() {
			return v
		}
		if e, ok := v.Interface().(embedded); ok {
			return reflect.ValueOf(e.S)
		}
		return v
	}

	hookForPtr := func(v reflect.Value) reflect.Value {
		if !v.CanInterface() {
			return v
		}
		if e, ok := v.Interface().(embedded); ok {
			return reflect.ValueOf(e.V)
		}
		return v
	}

	cases := []struct {
		Name   string
		Parts  []string
		Hook   ValueTransformationHookFn
		Input  interface{}
		Output interface{}
		Err    bool
	}{
		{
			Name:  "hook returns value of typed nil",
			Parts: []string{"Key"},
			Hook:  hookForPtr,
			Input: &struct {
				Key embedded
			}{Key: embedded{}},
			Output: (*embedded)(nil),
			Err:    false,
		},
		{
			Name:  "hook returns value of nil interface",
			Parts: []string{"Key"},
			Hook:  hookForInterface,
			Input: &struct {
				Key embedded
			}{Key: embedded{}},
			Output: nil,
			Err:    true,
		},
		{
			Name:   "top level don't replace",
			Parts:  []string{},
			Hook:   hookForInterface,
			Input:  embedded{S: "foo"},
			Output: embedded{S: "foo"},
			Err:    false,
		},
		{
			Name:  "1 deep replace",
			Parts: []string{"Key"},
			Hook:  hookForInterface,
			Input: &struct {
				Key embedded
			}{Key: embedded{S: "foo"}},
			Output: "foo",
			Err:    false,
		},
		{
			Name:  "2 deep replace",
			Parts: []string{"Key", "S2"},
			Hook:  hookForInterface,
			Input: &struct {
				Key embedded2
			}{Key: embedded2{S2: embedded{S: "foo"}}},
			Output: "foo",
			Err:    false,
		},
		{
			Name:  "two levels not last",
			Parts: []string{"Key"},
			Hook:  hookForInterface,
			Input: &struct {
				Key embedded
			}{Key: embedded{S: embedded2{S2: "foo"}}},
			Output: embedded2{S2: "foo"},
			Err:    false,
		},
		{
			Name:  "dont call hook twice per part",
			Parts: []string{"Key"},
			Hook:  hookForInterface,
			Input: &struct {
				Key embedded
			}{Key: embedded{S: embedded{S: "foo"}}},
			Output: embedded{S: "foo"},
			Err:    false,
		},
		{
			Name:  "through map",
			Parts: []string{"Key"},
			Hook:  hookForInterface,
			Input: &map[string]interface{}{
				"Key": embedded{S: "foo"},
			},
			Output: "foo",
			Err:    false,
		},
		{
			Name:   "slice key",
			Parts:  []string{"0"},
			Hook:   hookForInterface,
			Input:  []interface{}{embedded{S: "foo"}},
			Output: "foo",
			Err:    false,
		},
		{
			Name:   "array key",
			Parts:  []string{"0"},
			Hook:   hookForInterface,
			Input:  [1]interface{}{embedded{S: "foo"}},
			Output: "foo",
			Err:    false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			p := &Pointer{
				Parts: tc.Parts,
				Config: Config{
					ValueTransformationHook: tc.Hook,
				},
			}
			actual, err := p.Get(tc.Input)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s", err)
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(actual, tc.Output) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}

func TestPointerGet(t *testing.T) {
	type testStringType string
	type testIntType int

	cases := []struct {
		Name    string
		Parts   []string
		TagName string
		Input   interface{}
		Output  interface{}
		Err     bool
	}{
		{
			"empty",
			[]string{},
			"",
			42,
			42,
			false,
		},

		{
			"nil",
			nil,
			"",
			42,
			42,
			false,
		},

		{
			"map key",
			[]string{"foo"},
			"",
			map[string]interface{}{"foo": "bar"},
			"bar",
			false,
		},

		{
			"map key type change",
			[]string{"foo"},
			"",
			map[testStringType]interface{}{"foo": "bar"},
			"bar",
			false,
		},

		{
			"map key type change non-string",
			[]string{"42"},
			"",
			map[testIntType]interface{}{42: "bar"},
			"bar",
			false,
		},

		{
			"map key missing",
			[]string{"foo"},
			"",
			map[string]interface{}{"bar": "baz"},
			nil,
			true,
		},

		{
			"map key number",
			[]string{"42"},
			"",
			map[int]interface{}{42: "baz"},
			"baz",
			false,
		},

		{
			"map recursive",
			[]string{"foo", "42"},
			"",
			map[string]interface{}{
				"foo": map[int]interface{}{
					42: "baz",
				},
			},
			"baz",
			false,
		},

		{
			"slice key",
			[]string{"3"},
			"",
			[]interface{}{"a", "b", "c", "d", "e"},
			"d",
			false,
		},

		{
			"slice key non-existent",
			[]string{"7"},
			"",
			[]interface{}{"a", "b", "c", "d", "e"},
			nil,
			true,
		},

		{
			"slice key below zero",
			[]string{"-1"},
			"",
			[]interface{}{"a", "b", "c", "d", "e"},
			nil,
			true,
		},

		{
			"array key",
			[]string{"3"},
			"",
			&[5]interface{}{"a", "b", "c", "d", "e"},
			"d",
			false,
		},

		{
			"struct key",
			[]string{"Key"},
			"",
			&struct{ Key string }{Key: "foo"},
			"foo",
			false,
		},

		{
			"struct tag",
			[]string{"synthetic-name"},
			"",
			&struct {
				Key string `pointer:"synthetic-name"`
			}{Key: "foo"},
			"foo",
			false,
		},

		{
			"struct tag alt name",
			[]string{"synthetic-name"},
			"altptr",
			&struct {
				Key   string `altptr:"synthetic-name"`
				Other string `pointer:"synthetic-name"`
			}{Key: "foo", Other: "bar"},
			"foo",
			false,
		},

		{
			"struct tag ignore",
			[]string{"Key"},
			"altptr",
			&struct {
				Key string `altptr:"-"`
			}{Key: "foo"},
			"",
			true,
		},

		{
			"struct tag ignore and override",
			[]string{"X"},
			"",
			&struct {
				X string `pointer:"-"`
				Y string `pointer:"X"`
			}{X: "foo", Y: "bar"},
			"bar",
			false,
		},

		{
			"struct tag ignore after comma",
			[]string{"synthetic"},
			"pointer",
			&struct {
				Key string `pointer:"synthetic,name"`
			}{Key: "foo"},
			"foo",
			false,
		},

		{
			"struct tag invalid",
			[]string{"synthetic|name"},
			"",
			&struct {
				Key string `pointer:"synthetic|name"`
			}{Key: "foo"},
			"",
			true,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			p := &Pointer{Parts: tc.Parts, Config: Config{TagName: tc.TagName}}
			actual, err := p.Get(tc.Input)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s", err)
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(actual, tc.Output) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}
