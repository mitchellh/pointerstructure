package pointerstructure

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPointerParent(t *testing.T) {
	cases := []struct {
		Name     string
		Input    []string
		Expected []string
	}{
		{
			"basic",
			[]string{"foo", "bar"},
			[]string{"foo"},
		},

		{
			"single element",
			[]string{"foo"},
			[]string{},
		},

		{
			"root",
			nil,
			nil,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			p := &Pointer{Parts: tc.Input}
			result := p.Parent()

			actual := result.Parts
			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}
