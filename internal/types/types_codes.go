package types

import (
	"fmt"
	"net/http"
)

type (
	ResultCode   uint32
	ResultStatus uint32
)

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
	ErrRecoverPanic
	ErrCantGetProxyObject
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
		ErrUnknownTypeObjectPeer:  "unknown type obejct peer",
		ErrCantGetProxyObject:     "can't get proxy object",
		ErrCantInitProxy:          "can't init proxy",
		ErrPeerIsFailed:           "peer is failed",
		ErrRecoverPanic:           "recover panic",
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
