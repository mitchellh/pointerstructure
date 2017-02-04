package pointerstructure

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected []string
		Err      bool
	}{
		{
			"empty",
			"",
			nil,
			false,
		},

		{
			"relative",
			"foo",
			nil,
			true,
		},

		{
			"basic",
			"/foo/bar",
			[]string{"foo", "bar"},
			false,
		},

		{
			"three parts",
			"/foo/bar/baz",
			[]string{"foo", "bar", "baz"},
			false,
		},

		{
			"escaped /",
			"/foo/a~1b",
			[]string{"foo", "a/b"},
			false,
		},

		{
			"escaped ~",
			"/foo/a~0b",
			[]string{"foo", "a~b"},
			false,
		},

		{
			"escaped ~1",
			"/foo/a~01b",
			[]string{"foo", "a~1b"},
			false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			p, err := Parse(tc.Input)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s", err)
			}
			if err != nil {
				return
			}

			if p == nil {
				t.Fatal("nil pointer")
			}

			if !reflect.DeepEqual(p.Parts, tc.Expected) {
				t.Fatalf("bad: %#v", p.Parts)
			}
		})
	}
}
