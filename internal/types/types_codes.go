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

var ErrFileIsNil = errors.New("file is nil")
var ErrMaxCountAttempts = errors.New("exceeded maximum number of attempts")
var ErrRecoveredPanic = errors.New("recovered from panic")
var ErrInitProxy = errors.New("can't init proxy")
var ErrConfigPeersIsNil = errors.New("empty list peer in config file") 
var ErrCantGetProxyPort = errors.New("can't get proxy port")
var ErrCantGetProxySection = errors.New("can't get proxy section in config")

const (
	ResultOK ResultCode = iota
	ErrProxy
	ErrClientProxy
	ErrEmptyValue
	ErrCantFindFile
	ErrCantFindActivePeers
	ErrUnknownTypeObjectPeer
	ErrPeerIsFailed
	ErrCantMarshalJSON
	ErrSendResponseToClient
	ErrCopyDataPeerToClient
	ErrCopyDataClientToPeer
	ErrCountAttempts
	ErrGotPanic
	ErrCantInitProxy
	StatusBadRequest          ResultCode = http.StatusBadRequest
	StatusInternalServerError ResultCode = http.StatusInternalServerError
	StatusNotExtended         ResultCode = http.StatusNotExtended
	ResultUnknown             ResultCode = 0xFFFFFFFF
)

// ToStr returns a string representation of the ResultCode.
func (s ResultCode) ToStr() string {
	mapStatus := map[ResultCode]string{
		ResultOK:                  "SUCCESS",
		ErrProxy:                  "proxy error",
		ErrClientProxy:            "proxy client error",
		ErrCantMarshalJSON:        "can't marshal json object",
		ErrCopyDataPeerToClient:   "error copy data from peer to client",
		ErrCopyDataClientToPeer:   "error copy data from client to peer",
		ErrCountAttempts:          "exceeded maximum number of attempts",
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

// ToUint converts the ResultCode to its uint32 representation.
func (s ResultCode) ToUint() uint32 {
	return uint32(s)
}
