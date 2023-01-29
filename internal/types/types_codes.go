package types

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	ResultCode   uint32
	ResultStatus uint32
)

const ErrorCopyDataPeerToClient = "error copying data from peer to client, %v\n"
const ErrorCopyDataClientToPeer = "error copying data from client to peer, %v\n"

var ErrFileIsNil = errors.New("file is nil")

const (
	ResultOK ResultCode = iota
	ErrProxy
	ErrClientProxy
	ErrEmptyValue
	ErrCantFindFile
	ErrCantFindActivePeers
	ErrPeerIsFailed
	ErrCantMarshalJSON
	ErrSendResponseToClient
	ErrCopyDataPeerToClient
	ErrCopyDataClientToPeer
	ErrGotPanic
	StatusBadRequest          ResultCode = http.StatusBadRequest
	StatusInternalServerError ResultCode = http.StatusInternalServerError
	StatusNotExtended         ResultCode = http.StatusNotExtended
	ResultUnknown             ResultCode = 0xFFFFFFFF
)

func (s ResultCode) ToStr() string {
	mapStatus := map[ResultCode]string{
		ResultOK:                  "SUCCESS",
		ErrProxy:                  "proxy error",
		ErrClientProxy:            "proxy client error",
		ErrCantMarshalJSON:        "can't marshal json object",
		ErrCopyDataPeerToClient:   "error copy data from peer to client",
		ErrCopyDataClientToPeer:   "error copy data from client to peer",
		ErrSendResponseToClient:   "proxy err send response to client",
		StatusBadRequest:          "status bad request",
		StatusInternalServerError: "status internal server error",
		StatusNotExtended:         "status not extended",
		ErrEmptyValue:             "value is empty",
		ErrCantFindFile:           "can't find file",
		ErrCantFindActivePeers:    "can't find active peers",
		ErrGotPanic:               "got panic",
		ResultUnknown:             "unknown error",
	}

	m, ok := mapStatus[s]
	if !ok {
		return fmt.Sprintf("unknown result code: %d", s)
	}

	return m
}

func (s ResultCode) ToUint() uint32 {
	return uint32(s)
}