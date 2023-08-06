package core_test

import (
	"testing"
	"vbalancer/internal/core"
)

// TestTrimLastChar tests the TrimLastChar function.
func TestTrimLastChar(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input string
		want  string
	}{
		{"value", "valu"},
		{"", ""},
		{"abc", "ab"},
		{"hello world", "hello worl"},
		{"世界", "世"},
		{"こんにちは", "こんにち"},
		{"你好", "你"},
		{"Привет", "Приве"},
	}
	for _, tc := range testCases {
		got := core.TrimLastChar(tc.input)
		if got != tc.want {
			t.Errorf("TrimLastChar(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// TestFmtStringWithDelimiter tests the FmtStringWithDelimiter function.
func TestFmtStringWithDelimiter(t *testing.T) {
	t.Parallel()

	delimiter := ","
	tests := []struct {
		values  []interface{}
		result  string
		isError bool
	}{
		{[]interface{}{"a", "b", "c"}, "a,b,c", false},
		{[]interface{}{1, 2, 3}, "1,2,3", false},
		{[]interface{}{"a", 1, []int{1, 2, 3}}, "a,1,1,2,3", false},
		{[]interface{}{"a", nil, "c"}, "a,c", false},
		{[]interface{}{"a", "b", "c\nd"}, "a,b,c d", false},
		{[]interface{}{"a", "b", "c", []int{1, 2, 3}}, "a,b,c,1,2,3", false},
		{[]interface{}{1, 2, []string{"a", "b", "c"}}, "1,2,a,b,c", false},
		{[]interface{}{1, 2, "c", []int{}}, "1,2,c", false},
	}

	for _, test := range tests {
		result := core.FmtStringWithDelimiter(delimiter, test.values...)
		if result != test.result {
			t.Errorf("Expected %s, but got %s", test.result, result)
		}
	}
}
