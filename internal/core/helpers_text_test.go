package core_test

import (
	"testing"
	"vbalancer/internal/core"
)


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