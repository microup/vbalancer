package response

import (
	"encoding/json"
	"fmt"
	"net"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

// Response is a struct that contains the response to the client.
type Response struct {
	StatusCode  types.ResultCode `json:"statusCode"`
	Description string           `json:"description"`
	logger      vlog.ILog
}

// New create a new response object.
func New(logger vlog.ILog) *Response {
	return &Response{
		StatusCode:  types.ResultUnknown,
		Description: "",
		logger:      logger,
	}
}

// SentResponseToClient send a response to the client.
func (r *Response) SentResponseToClient(client net.Conn, err error) error {
	r.StatusCode = types.ErrProxy
	r.Description = err.Error()

	responseJSON, err := json.Marshal(r)
	
	if err != nil {
		r.logger.Add(vlog.Debug,
			types.ErrCantMarshalJSON,
			vlog.RemoteAddr(client.RemoteAddr().String()),
			types.ErrCantMarshalJSON.ToStr())

		return fmt.Errorf("%w", err)
	}

	responseLen := len(responseJSON)
	responseBody :=
		fmt.Sprintf(
			"HTTP/1.1 200 OK\r\nConnection: close\r\n"+
			"Content-Type: application/json\r\n"+
			"Content-Length: %d\r\n\r\n"+
			"%s", responseLen, responseJSON)

		
	_, err = client.Write([]byte(responseBody))
	if err != nil {
		r.logger.Add(vlog.Debug, types.ErrSendResponseToClient, vlog.RemoteAddr(client.RemoteAddr().String()),
			types.ErrSendResponseToClient.ToStr())

		return fmt.Errorf("%w", err )
	}

	return nil
}

