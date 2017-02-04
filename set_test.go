package pointerstructure

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPointerSet(t *testing.T) {
	type testStringType string
	type testIntType int

	cases := []struct {
		Name   string
		Parts  []string
		Doc    interface{}
		Value  interface{}
		Output interface{}
		Err    bool
	}{
		{
			"empty",
			[]string{},
			42,
			[]string{},
			[]string{},
			false,
		},

		{
			"nil",
			nil,
			42,
			84,
			84,
			false,
		},

		{
			"map key",
			[]string{"foo"},
			map[string]interface{}{"foo": "bar"},
			"baz",
			map[string]interface{}{"foo": "baz"},
			false,
		},

		{
			"nested map key",
			[]string{"foo", "bar"},
			map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 42,
				},
			},
			"baz",
			map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
				},
			},
			false,
		},

		{
			"map value type coercse",
			[]string{"foo"},
			map[string]int{"foo": 12},
			"42",
			map[string]int{"foo": 42},
			false,
		},

		{
			"slice index",
			[]string{"0"},
			[]interface{}{42},
			"baz",
			[]interface{}{"baz"},
			false,
		},

		{
			"slice index non-zero",
			[]string{"1"},
			[]interface{}{42, 84, 168},
			"baz",
			[]interface{}{42, "baz", 168},
			false,
		},

		{
			"slice index append",
			[]string{"-"},
			[]interface{}{42},
			"baz",
			[]interface{}{42, "baz"},
			false,
		},

		{
			"slice index value coerce",
			[]string{"0"},
			[]int{42},
			"84",
			[]int{84},
			false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			p := &Pointer{Parts: tc.Parts}
			actual, err := p.Set(tc.Doc, tc.Value)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s", err)
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(actual, tc.Output) {
				t.Fatalf("bad: %#v != %#v", actual, tc.Output)
			}
		})
	}
}
