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
	ProxyError
	ErrEmptyValue
	ErrCantFindFile
	ErrCantFinePeers
	StatusBadRequest          ResultCode = http.StatusBadRequest
	StatusInternalServerError ResultCode = http.StatusInternalServerError
	StatusNotExtended         ResultCode = http.StatusNotExtended
	ErrUnknown                ResultCode = 0xFFFFFFFF
)

func (s ResultCode) ToStr() string {
	mapStatus := map[ResultCode]string{
		ResultOK:                  "SUCCESS",
		ProxyError:                "Proxy error",
		StatusBadRequest:          "StatusBadRequest",
		StatusInternalServerError: "StatusInternalServerError",
		StatusNotExtended:         "StatusNotExtended",
		ErrEmptyValue:             "Value is empty",
		ErrCantFindFile:           "Can't find file",
		ErrCantFinePeers:          "Can't find peers",
		ErrUnknown:                "unknown error",
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
