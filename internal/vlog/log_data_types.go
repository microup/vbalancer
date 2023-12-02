package vlog

import (
	"fmt"
	"strings"
	"time"
	"vbalancer/internal/core"
	"vbalancer/internal/types"
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
	TypeLog, types.ResultCode,
	RemoteAddr, PeerAddr, string) {
	var typeLog TypeLog
	var resultCode types.ResultCode //nolint:wsl
	var remoteAddr RemoteAddr       //nolint:wsl
	var peerAddr PeerAddr           //nolint:wsl
	var val strings.Builder         //nolint:wsl

	for _, value := range values {
		switch valueTypeLog := value.(type) {
		case TypeLog:
			typeLog = valueTypeLog
		case types.ResultCode:
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

// BuildRecord function takes in several input values such as log type, result code,
// remote address, client host, proxy host and string values. It creates a record in the
// specified format with the current date and time by using the "GetDateTimeStr" function
// from the "core" package. The result is returned as a tuple with the log type and the formatted record string.
func BuildRecord(
	typeLog TypeLog,
	resultCode types.ResultCode,
	remoteAddr RemoteAddr,
	peerAddr PeerAddr,
	valuesStr string) (TypeLog, string) {
	var recordTime = time.Now()

	var dateStr, timeStr = core.GetDateTimeStr(recordTime)

	var resultStr = core.FmtStringWithDelimiter(";", valuesStr)

	var recordRow = fmt.Sprintf("%s;%s;%s;%d;%s;%s;%s;",
		dateStr,
		timeStr,
		typeLog.GetStr(),
		resultCode,
		remoteAddr,
		peerAddr,
		resultStr)

	return typeLog, recordRow
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
