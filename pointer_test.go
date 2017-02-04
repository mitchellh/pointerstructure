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

func TestPointerString(t *testing.T) {
	cases := []string{
		"/foo",
		"/foo/bar",
		"/foo/bar~0",
		"/foo/bar~1",
		"/foo/bar~01/baz",
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			p := MustParse(tc)
			actual := p.String()
			if actual != tc {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}
