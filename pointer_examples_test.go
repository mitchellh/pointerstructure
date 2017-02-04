package pointerstructure

import (
	"fmt"
)

func ExampleGet() {
	complex := map[string]interface{}{
		"alice": 42,
		"bob": []interface{}{
			map[string]interface{}{
				"name": "Bob",
			},
		},
	}

	value, err := Get(complex, "/bob/0/name")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", value)
	// Output:
	// Bob
}

func ExampleSet() {
	complex := map[string]interface{}{
		"alice": 42,
		"bob": []interface{}{
			map[string]interface{}{
				"name": "Bob",
			},
		},
	}

	value, err := Set(complex, "/bob/0/name", "Alice")
	if err != nil {
		panic(err)
	}

	value, err = Get(complex, "/bob/0/name")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", value)
	// Output:
	// Alice
}
