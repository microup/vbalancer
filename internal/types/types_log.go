package types

import (
	"fmt"
	"strings"
)

type (
	TypeLog uint8
	IsSave  bool
)

type (
	RemoteAddr      string
	PeerAddr        string
	ProxyRequestURI string
)

const (
	Disable TypeLog = iota
	Fatal
	Error
	Warning
	Debug
	Info
)

// ParseValues takes in a slice of interface values and returns multiple values of different types.
// It loops through the values in the input slice and assigns each value to its corresponding variable
// based on its type. The function returns TypeLog, ResultCode, RemoteAddr, ClientHost, ProxyHost
// and a string value built
// from concatenating the string values in the input slice with semicolons.
func ParseValues(values []interface{}) (
	TypeLog, ResultCode,
	RemoteAddr, PeerAddr, string) {
	var typeLog TypeLog
	var resultCode ResultCode
	var remoteAddr RemoteAddr
	var peerAddr PeerAddr
	var val strings.Builder

	for _, value := range values {
		switch valueTypeLog := value.(type) {
		case TypeLog:
			typeLog = valueTypeLog
		case ResultCode:
			resultCode = valueTypeLog
		case string:
			val.WriteString(valueTypeLog)
			val.WriteString(";")
		case error:
			val.WriteString(valueTypeLog.Error())
			val.WriteString(";")
		case RemoteAddr:
			remoteAddr = valueTypeLog
		case PeerAddr:
			peerAddr = valueTypeLog
		}
	}

	return typeLog, resultCode, remoteAddr, peerAddr, val.String()
}

// GetStr this function returns a string representation of a "TypeLog"
// value by mapping it to its string equivalent. If the "TypeLog" value is
// not found in the map, it returns an error message indicating an unknown result code.
func (s TypeLog) GetStr() string {
	mapTypeLog := map[TypeLog]string{
		Disable: "DISABLE",
		Info:    "INFO",
		Debug:   "DEBUG",
		Warning: "WARNING",
		Error:   "ERROR",
		Fatal:   "FATAL",
	}

	m, ok := mapTypeLog[s]
	if !ok {
		return fmt.Sprintf("unknown result code: %d", s)
	}

	return m
}
