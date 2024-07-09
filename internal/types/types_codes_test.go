package types_test

import (
	"testing"
	"vbalancer/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestResultCodeToStr(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		input types.ResultCode
		want  string
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
		{types.ErrRecoverPanic, "recover panic"},
		{types.ResultUnknown, "unknown error"},
		{types.ResultCode(0xABC), "unknown result code: 2748"},
	}

	for _, test := range testCases {
		assert.Equalf(t, test.want, test.input.ToStr(), "input: `%d | want %d`", test.input, test.want)
	}
}

func TestResultCodeToUint(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		input types.ResultCode
		want  uint32
	}{
		{types.ResultOK, 0},
		{types.ResultUnknown, 4294967295},
		{types.ResultCode(0xABC), 2748},
	}

	for _, test := range testCases {
		assert.Equalf(t, test.want, test.input.ToUint(), "input: `%d | want %d`", test.input, test.want)
	}
}
