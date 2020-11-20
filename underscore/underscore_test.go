package underscore

import "testing"

func TestCamel(t *testing.T) {
	testCases := map[string]struct {
		arg  string
		want string
	}{
		"some_name":    {"thisIsACamelCaseString", "this_is_a_camel_case_string"},
		"another_name": {"with a space", "with a space"},
		"with a space": {"endsWithA", "ends_with_a"},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Logf("Testing: %q", tt.arg)
			if got := Camel(tt.arg); got != tt.want {
				// Fatalf breaks this subtest but not all tests
				t.Fatalf("Camel(%q) = %q; want %q", tt.arg, got, tt.want)
			}
		})
	}
}
