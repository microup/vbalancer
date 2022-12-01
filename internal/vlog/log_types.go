package vlog

import (
	"fmt"
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
