package vlog

import (
	"fmt"
	"strings"
	"time"
	"vbalancer/internal/core"
	"vbalancer/internal/types"
)

// Package file is a simple file logger.
type (
	// TypeLog is a type of log level.
	TypeLog uint8
	// IsSave is a type of bool.
	IsSave bool
)

type (
	// RemoteAddr is a type of string.
	RemoteAddr string
	// ClientHost is a type of string.
	ClientHost string
	// ClientMethod is a type of string.
	ClientMethod string
	// ClientProto is a type of string.
	ClientProto string
	// ClientURI is a type of string.
	ClientURI string
	// ProxyHost is a type of string.
	ProxyHost string
	// ProxyMethod is a type of string.
	ProxyMethod string
	// ProxyProto is a type of string.
	ProxyProto string
	// ProxyURI is a type of string.
	ProxyURI string
	// ProxyRequestURI is a type of string.
	ProxyRequestURI string
)

const (
	// Disable is a type of TypeLog.
	Disable TypeLog = iota
	// Fatal is a type of TypeLog.
	Fatal
	// Error is a type of TypeLog.
	Error
	// Warning is a type of TypeLog.
	Warning
	// Debug is a type of TypeLog.
	Debug
	// Info is a type of TypeLog.
	Info
)

// ParseValues - takes in a slice of interface values and returns multiple values of different types.
// It loops through the values in the input slice and assigns each value to its corresponding variable
// based on its type. The function returns TypeLog, ResultCode, RemoteAddr, ClientHost, ClientMethod,
// ClientProto, ClientURI, ProxyHost, ProxyMethod, ProxyProto, ProxyURI and a string value built
// from concatenating the string values in the input slice with semicolons.
//
//nolint:cyclop,funlen
func ParseValues(values []interface{}) (
	TypeLog, types.ResultCode,
	RemoteAddr, ClientHost, ClientMethod,
	ClientProto, ClientURI, ProxyHost,
	ProxyMethod, ProxyProto, ProxyURI, string) {
	var typeLog TypeLog
	var resultCode types.ResultCode //nolint:wsl
	var remoteAddr RemoteAddr       //nolint:wsl
	var clientHost ClientHost       //nolint:wsl
	var clientMethod ClientMethod   //nolint:wsl
	var clientProto ClientProto     //nolint:wsl
	var clientURI ClientURI         //nolint:wsl
	var proxyHost ProxyHost         //nolint:wsl
	var proxyMethod ProxyMethod     //nolint:wsl
	var proxyProto ProxyProto       //nolint:wsl
	var proxyURI ProxyURI           //nolint:wsl
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
		case ClientHost:
			clientHost = valueTypeLog
		case ClientMethod:
			clientMethod = valueTypeLog
		case ClientProto:
			clientProto = valueTypeLog
		case ClientURI:
			clientURI = valueTypeLog
		case ProxyHost:
			proxyHost = valueTypeLog
		case ProxyMethod:
			proxyMethod = valueTypeLog
		case ProxyProto:
			proxyProto = valueTypeLog
		case ProxyURI:
			proxyURI = valueTypeLog
		}
	}

	return typeLog, resultCode, remoteAddr,
		clientHost, clientMethod, clientProto, clientURI,
		proxyHost, proxyMethod, proxyProto, proxyURI, val.String()
}

// BuildRecord - function takes in several input values such as log type, result code,
// remote address, client host, client method, client protocol, client URI, proxy host,
// proxy method, proxy protocol, proxy URI, and string values. It creates a record in the
// specified format with the current date and time by using the "GetDateTimeStr" function
// from the "core" package. The result is returned as a tuple with the log type and the formatted record string.
func BuildRecord(typeLog TypeLog, resultCode types.ResultCode,
	remoteAddr RemoteAddr, clientHost ClientHost, clientMethod ClientMethod,
	clientProto ClientProto, clientURI ClientURI, proxyHost ProxyHost,
	proxyMethod ProxyMethod, proxyProto ProxyProto, proxyURI ProxyURI,
	valuesStr string) (TypeLog, string) {
	var recordTime = time.Now()

	var dateStr, timeStr = core.GetDateTimeStr(recordTime)

	var resultFmtStr = core.FmtStringWithDelimiter(";", valuesStr)

	var recordRow = fmt.Sprintf("%s;%s;%s;%d;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s",
		dateStr,
		timeStr,
		typeLog.GetStr(),
		resultCode,
		remoteAddr,
		clientHost,
		clientMethod,
		clientProto,
		clientURI,
		proxyMethod,
		proxyProto,
		proxyHost,
		proxyURI,
		resultFmtStr)

	return typeLog, recordRow
}

// GetStr - this function returns a string representation of a "TypeLog"
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
