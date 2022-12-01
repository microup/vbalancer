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
	ResultOK                  ResultCode = 0
	ProxyError                           = 1
	StatusBadRequest                     = http.StatusBadRequest
	StatusInternalServerError            = http.StatusInternalServerError
	StatusNotExtended                    = http.StatusNotExtended
	ErrEmptyValue                        = 1001
	ErrCantFinePeers                     = 2001
	ErrUnknown                           = 0xFFFFFFFF
)

func (s ResultCode) ToStr() string {
	mapStatus := map[ResultCode]string{
		ResultOK:                  "SUCCESS",
		ProxyError:                "Proxy error",
		StatusBadRequest:          "StatusBadRequest",
		StatusInternalServerError: "StatusInternalServerError",
		StatusNotExtended:         "StatusNotExtended",
		ErrEmptyValue:             "Value is empty",
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
