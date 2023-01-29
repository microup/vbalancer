package types

import (
	"fmt"
	"strings"
	"time"
	"vbalancer/internal/core"
)

type (
	TypeLog uint8
	IsSave  bool
)

type (
	RemoteAddr      string
	ClientHost      string
	ClientMethod    string
	ClientProto     string
	ClientURI       string
	ProxyHost       string
	ProxyMethod     string
	ProxyProto      string
	ProxyURI        string
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

//nolint:cyclop
func ParseValues(values []interface{}) (
	TypeLog, ResultCode,
	RemoteAddr, ClientHost, ClientMethod,
	ClientProto, ClientURI, ProxyHost,
	ProxyMethod, ProxyProto, ProxyURI, string) {
	var typeLog TypeLog
	var resultCode ResultCode     //nolint:wsl
	var remoteAddr RemoteAddr     //nolint:wsl
	var clientHost ClientHost     //nolint:wsl
	var clientMethod ClientMethod //nolint:wsl
	var clientProto ClientProto   //nolint:wsl
	var clientURI ClientURI       //nolint:wsl
	var proxyHost ProxyHost       //nolint:wsl
	var proxyMethod ProxyMethod   //nolint:wsl
	var proxyProto ProxyProto     //nolint:wsl
	var proxyURI ProxyURI         //nolint:wsl
	var val strings.Builder       //nolint:wsl

	for _, value := range values {
		switch valueTypeLog := value.(type) {
		case TypeLog:
			typeLog = valueTypeLog
		case ResultCode:
			resultCode = valueTypeLog
		case string:
			val.WriteString(valueTypeLog)
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

func BuildRecord(typeLog TypeLog, resultCode ResultCode,
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
		*resultFmtStr)

	return typeLog, recordRow
}
