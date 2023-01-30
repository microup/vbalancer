package vlog_test

import (
	"sync"
	"testing"
	"vbalancer/internal/vlog"
)

func TestGetRecordsText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputVLog           vlog.VLog
		inputIsReverse      bool
		expectedOutput      string
		expectedErrorOutput error
	}{
		{
			inputVLog: vlog.VLog{ //nolint:exhaustivestruct,exhaustruct
				Mu:                &sync.Mutex{},
				MapLastLogRecords: []string{"Record 1", "Record 2", "Record 3"},
			},
			inputIsReverse:      false,
			expectedOutput:      "Record 1<BR>Record 2<BR>Record 3",
			expectedErrorOutput: nil,
		},
		{
			inputVLog: vlog.VLog{ //nolint:exhaustivestruct,exhaustruct
				Mu:                &sync.Mutex{},
				MapLastLogRecords: []string{"Record 1", "Record 2", "Record 3"},
			},
			inputIsReverse:      true,
			expectedOutput:      "Record 3<BR>Record 2<BR>Record 1",
			expectedErrorOutput: nil,
		},
		{
			inputVLog: vlog.VLog{ //nolint:exhaustivestruct,exhaustruct
				Mu: &sync.Mutex{},
			},
			inputIsReverse:      false,
			expectedOutput:      "",
			expectedErrorOutput: nil,
		},
		{
			inputVLog: vlog.VLog{ //nolint:exhaustivestruct,exhaustruct
				Mu:                &sync.Mutex{},
				MapLastLogRecords: []string{},
			},
			inputIsReverse:      false,
			expectedOutput:      "",
			expectedErrorOutput: nil,
		},
	}

	for _, test := range tests {
		output := test.inputVLog.GetRecordsText(test.inputIsReverse)

		if output != test.expectedOutput {
			t.Errorf("Test case %v failed: expected %v, got %v\n", test, test.expectedOutput, output)
		}
	}
}
