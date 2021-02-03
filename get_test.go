package pointerstructure

import (
	"fmt"
	"reflect"
	"testing"
)

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
