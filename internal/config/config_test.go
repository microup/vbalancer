package config_test

import (
	"fmt"
	"os"
	"testing"
	"vbalancer/internal/config"
	"vbalancer/internal/types"
)

//nolint:funlen
func TestInitProxyPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		envVar     string
		expected   types.ResultCode
		checkValue string
	}{
		{
			name:       "invalid proxy port from env var",
			envVar:     ":",
			expected:   types.ErrEmptyValue,
			checkValue: "",
		},
		{
			name:       "invalid proxy port from env var - 2",
			envVar:     "",
			expected:   types.ResultOK,
			checkValue: fmt.Sprintf(":%d", config.DefaultProxyPort),
		},
		{
			name:       "empty proxy port from env var",
			envVar:     "",
			expected:   types.ResultOK,
			checkValue: fmt.Sprintf(":%d", config.DefaultProxyPort),
		},
		{
			name:       "valid proxy port from env var",
			envVar:     "8080",
			expected:   types.ResultOK,
			checkValue: ":8080",
		},
		{
			name:       "empty proxy port from default value",
			envVar:     " ",
			expected:   types.ErrEmptyValue,
			checkValue: ":",
		},
		{
			name:       "empty proxy port from default value",
			envVar:     "          ",
			expected:   types.ErrEmptyValue,
			checkValue: ":",
		},
	}

	config := &config.Config{
		Logger:    nil,
		Proxy:     nil,
		Peers:     nil,
		ProxyPort: "",
	}

	for _, test := range tests {
		config.ProxyPort = ""

		os.Clearenv()
		os.Setenv("ProxyPort", test.envVar)

		result := config.InitProxyPort()

		if result != test.expected {
			t.Fatalf("name: %s, expected result %v, got %v", test.name, test.expected, result)
		}

		if config.ProxyPort != test.checkValue {
			t.Fatalf("name: %s, expected value %s, got %s", test.name, test.checkValue, config.ProxyPort)
		}
	}
}
