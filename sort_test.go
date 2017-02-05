package pointerstructure

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestSort_interface(t *testing.T) {
	var _ sort.Interface = new(PointerSlice)
}

func TestSort(t *testing.T) {
	cases := []struct {
		Input  []string
		Output []string
	}{
		{
			[]string{"/foo", ""},
			[]string{"", "/foo"},
		},

		{
			[]string{"/foo", "", "/foo/0"},
			[]string{"", "/foo", "/foo/0"},
		},

		{
			[]string{"/foo", "", "/bar/0"},
			[]string{"", "/bar/0", "/foo"},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var ps []*Pointer
			for _, raw := range tc.Input {
				ps = append(ps, MustParse(raw))
			}

			Sort(ps)

			result := make([]string, len(ps))
			for i, p := range ps {
				result[i] = p.String()
			}

			if !reflect.DeepEqual(result, tc.Output) {
				t.Fatalf("bad: %#v", result)
			}
		})
	}
}
