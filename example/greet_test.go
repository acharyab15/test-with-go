package example

import "fmt"

func ExampleDemo_Hello() {
	greeting, err := Hello("Bipeen")
	if err != nil {
		panic(err)
	}
	fmt.Println(greeting)

	// Output: Hello, Bipeen
}

func ExamplePage() {
	checkIns := map[string]bool{
		"Bob":   true,
		"Alice": false,
		"Eve":   false,
		"Janet": true,
		"Susan": true,
		"John":  false,
	}
	Page(checkIns)
	// Unordered Output:
	// Paging Alice; please see the front desk to check in.
	// Paging Eve; please see the front desk to check in.
	// Paging John; please see the front desk to check in.
}
