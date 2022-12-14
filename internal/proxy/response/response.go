package response

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

type Response struct {
	StatusCode  types.ResultCode `json:"statusCode"`
	Description string           `json:"description"`
	logger      *vlog.VLog
}

func New(logger *vlog.VLog) *Response {
	return &Response{
		StatusCode:  types.ResultUnknown,
		Description: "",
		logger:      logger,
	}
}

func (r *Response) SentResponse(client net.Conn, codeResponse types.ResultCode) {
	r.StatusCode = codeResponse
	r.Description = codeResponse.ToStr()

	responseJSON, err := json.Marshal(r)

	if err != nil {
		r.logger.Add(vlog.Debug, types.ErrCantMarshalJSON, vlog.RemoteAddr(client.RemoteAddr().String()),
			types.ErrCantMarshalJSON.ToStr())

		return
	}

	responseLen := len(responseJSON)
	responseBody :=
		fmt.Sprintf("HTTP/1.1 200 OK\r\nConnection: close\r\n"+
			"Content-Type: application/json\r\n"+
			"Content-Length: %d\r\n\r\n"+
			"%s", responseLen, responseJSON)

	_, err = client.Write([]byte(responseBody))
	if err != nil {
		r.logger.Add(vlog.Debug, types.ErrSendResponseToClient, vlog.RemoteAddr(client.RemoteAddr().String()),
			types.ErrSendResponseToClient.ToStr())
	}
	
	time.Sleep(1 * time.Nanosecond)
}
