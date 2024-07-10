//nolint:testpackage // the need to test the private method updatePort() in the proxy struct.
package proxy

import (
	"fmt"
	"testing"

	"vbalancer/internal/types"
	"github.com/stretchr/testify/assert"

)

// TestGetProxyPort tests the UpdatePort function.
// It validates UpdatePort handles invalid environment variable values,
// default values, and valid custom environment variable values correctly.
func TestGetProxyPort(t *testing.T) {
	testCases := []struct {
		port      string
		name      string
		envVar    string
		want      types.ResultCode
		wantValue string
	}{
		{
			name:      "set port `1234`",
			port:      "1234",
			envVar:    ":",
			want:      types.ResultOK,
			wantValue: ":1234",
		},
		{
			name:      "empty env var, got DefaultPort",
			port:      "",
			envVar:    "",
			want:      types.ResultOK,
			wantValue: fmt.Sprintf(":%s", types.DefaultProxyPort),
		},
		{
			name:      "valid proxy port from env var",
			port:      "",
			envVar:    "8080",
			want:      types.ResultOK,
			wantValue: ":8080",
		},
		{
			name:      "empty proxy port from default value",
			envVar:    " ",
			want:      types.ErrEmptyValue,
			wantValue: ":",
			port:      ":",
		},
		{
			name:      "empty proxy port from default value",
			envVar:    "          ",
			want:      types.ErrEmptyValue,
			wantValue: ":",
			port:      ":",
		},
	}

	prx := &Proxy{
		Logger:                nil,
		Port:                  "",
		ClientDeadLineTime:    10,
		PeerConnectionTimeout: 10,
		MaxCountConnection:    100,
		Peers:                 nil,
		Rules:                 nil,
		notify: 				make(chan error),
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			prx.Port = testCase.port

			t.Setenv(types.ProxyPort, testCase.envVar)

			result := prx.updatePort()

			assert.Equal(t, testCase.want, result, "name: `%s`")

			assert.Equal(t, testCase.wantValue, prx.Port, "name: `%s`")
		})
	}
}
