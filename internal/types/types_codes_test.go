package types_test

import (
	"testing"
	"vbalancer/internal/types"
)

func TestResultCodeToStr(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input    types.ResultCode
		expected string
	}{
		{types.ResultOK, "SUCCESS"},
		{types.ErrProxy, "proxy error"},
		{types.ErrClientProxy, "proxy client error"},
		{types.ErrCantMarshalJSON, "can't marshal json object"},
		{types.ErrCopyDataPeerToClient, "error copy data from peer to client"},
		{types.ErrCopyDataClientToPeer, "error copy data from client to peer"},
		{types.ErrSendResponseToClient, "proxy err send response to client"},
		{types.StatusBadRequest, "status bad request"},
		{types.StatusInternalServerError, "status internal server error"},
		{types.StatusNotExtended, "status not extended"},
		{types.ErrEmptyValue, "value is empty"},
		{types.ErrCantFindFile, "can't find file"},
		{types.ErrCantFindActivePeers, "can't find active peers"},
		{types.ErrGotPanic, "got panic"},
		{types.ResultUnknown, "unknown error"},
		{types.ResultCode(0xABC), "unknown result code: 2748"},
	}

	for _, test := range tests {
		if result := test.input.ToStr(); result != test.expected {
			t.Errorf("TestResultCodeToStr(%d) = %s; expected = %s", test.input, result, test.expected)
		}
	}
}

func TestResultCodeToUint(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input    types.ResultCode
		expected uint32
	}{
		{types.ResultOK, 0},
		{types.ErrProxy, 1},
		{types.ErrClientProxy, 2},
		{types.ErrCantMarshalJSON, 7},
		{types.ErrCopyDataPeerToClient, 9},
		{types.ErrCopyDataClientToPeer, 10},
		{types.ErrSendResponseToClient, 8},
		{types.StatusBadRequest, 400},
		{types.StatusInternalServerError, 500},
		{types.StatusNotExtended, 510},
		{types.ErrEmptyValue, 3},
		{types.ErrCantFindFile, 4},
		{types.ErrCantFindActivePeers, 5},
		{types.ErrGotPanic, 12},
		{types.ResultUnknown, 4294967295},
		{types.ResultCode(0xABC), 2748},
	}

	for _, test := range tests {
		if result := test.input.ToUint(); result != test.expected {
			t.Errorf("TestResultCodeToUint(%d) = %d; expected = %d", test.input, result, test.expected)
		}
	}
}
