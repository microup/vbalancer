package vlog

import (
	"fmt"

	"log"
	"sync"
	"time"

	"os"

	"vbalancer/internal/core"
	"vbalancer/internal/types"
)

type VLog struct {
	cfg               *Config
	fileLog           *os.File
	countToLogID      int
	mapLastLogRecords []string
	mu                sync.Mutex
	headerCSV         string
	startTimeLog      time.Time
	wgNewLog          sync.WaitGroup
	IsDisabled        bool
}

func New(cfg *Config) (*VLog, error) {

	l := &VLog{
		cfg:          cfg,
		countToLogID: -1,
		headerCSV: fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;", "Date", "Time", "Type", "ResultCode", "RemoteAddr",
			"ClientHost", "ClientMethod", "ClientProto", "ClientURI", "PeerMethod", "PeerProto", "PeerHost",
			"PeerRequestURI", "Description"),
		startTimeLog: time.Now(),
		IsDisabled:   false,
	}

	err := l.newFileLog("", true)
	if err != nil {
		return nil, err
	}

	return l, nil

}

func (v *VLog) GetCountRecords() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.mapLastLogRecords == nil {
		return 0
	}
	return len(v.mapLastLogRecords)
}

func (v *VLog) Add(values ...interface{}) {
	if v.IsDisabled {
		return
	}
	go v.addInThread(values...)
}

func (v *VLog) addInThread(values ...interface{}) {
	v.wgNewLog.Wait()

	v.mu.Lock()
	defer v.mu.Unlock()

	if values == nil || v.mapLastLogRecords == nil {
		return
	}

	typeLog, recordRow := v.buildCsvRecord(values)

	switch typeLog {
	case Fatal:
		log.Fatal(recordRow)
	default:
		log.Print(recordRow)
	}

	_, err := v.fileLog.WriteString(recordRow + "\n")
	if err != nil {
		log.Printf("Error: %s - is writing: %s to log file: %s", err, recordRow, v.fileLog.Name())
	}

	v.removeOldRecordsFromMemory()

	v.mapLastLogRecords = append(v.mapLastLogRecords, recordRow)

	err = v.checkToCreateNewLogFile()
	if err != nil {
		log.Printf("Error: %s - is writing: %s to log file: %s", err, recordRow, v.fileLog.Name())
	}
}

func (v *VLog) buildCsvRecord(values []interface{}) (TypeLog, string) {
	var typeLog TypeLog

	var val string

	var resultCode types.ResultCode

	var remoteAddr string

	var clientHost string

	var clientMethod string

	var clientProto string

	var clientURI string

	var proxyHost string

	var proxyMethod string

	var proxyProto string

	var proxyURI string

	for _, value := range values {
		switch valueTypeLog := value.(type) {
		case TypeLog:
			typeLog = TypeLog(valueTypeLog)
		case types.ResultCode:
			resultCode = types.ResultCode(valueTypeLog)
		case string:
			val = val + string(valueTypeLog) + ","
		case RemoteAddr:
			remoteAddr = string(value.(RemoteAddr))
		case ClientHost:
			clientHost = string(value.(ClientHost))
		case ClientMethod:
			clientMethod = string(value.(ClientMethod))
		case ClientProto:
			clientProto = string(value.(ClientProto))
		case ClientURI:
			clientURI = string(value.(ClientURI))
		case ProxyHost:
			proxyHost = string(value.(ProxyHost))
		case ProxyMethod:
			proxyMethod = string(value.(ProxyMethod))
		case ProxyProto:
			proxyProto = string(value.(ProxyProto))
		case ProxyURI:
			proxyURI = string(value.(ProxyURI))
		}
	}

	resultFmtStr := core.FmtStringWithDelimiter(";", val)

	recordTime := time.Now()
	dateStr := recordTime.Format("2006-01-02")
	timeStr := recordTime.Format("15:04:05")

	recordRow := fmt.Sprintf("%s;%s;%s;%d;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s",
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
