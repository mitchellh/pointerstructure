package pointerstructure

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPointerDelete(t *testing.T) {
	type testStringType string
	type testIntType int

	cases := []struct {
		Name   string
		Parts  []string
		Doc    interface{}
		Output interface{}
		Err    bool
	}{
		{
			"empty",
			[]string{},
			42,
			nil,
			false,
		},

		{
			"nil",
			nil,
			42,
			nil,
			false,
		},

		{
			"map key",
			[]string{"foo"},
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{},
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
			map[string]interface{}{
				"foo": map[string]interface{}{},
			},
			false,
		},

		{
			"slice index",
			[]string{"0"},
			[]interface{}{42},
			[]interface{}{},
			false,
		},

		{
			"slice index non-zero",
			[]string{"1"},
			[]interface{}{42, 84, 168},
			[]interface{}{42, 168},
			false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			p := &Pointer{Parts: tc.Parts}
			actual, err := p.Delete(tc.Doc)
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
