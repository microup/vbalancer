package vlog

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"vbalancer/internal/core"
	"vbalancer/internal/types"
)

type VLog struct {
	cfg               *Config
	fileLog           *os.File
	countToLogID      int
	mapLastLogRecords []string
	mu                *sync.Mutex
	headerCSV         string
	startTimeLog      time.Time
	wgNewLog          *sync.WaitGroup
	IsDisabled        bool
}

func New(cfg *Config) (*VLog, error) {
	vLog := &VLog{
		wgNewLog:          &sync.WaitGroup{},
		mu:                &sync.Mutex{},
		fileLog:           nil,
		cfg:               cfg,
		countToLogID:      -1,
		mapLastLogRecords: []string{},
		headerCSV: fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;%s;",
			"Date", "Time", "Type", "ResultCode", "RemoteAddr",
			"ClientHost", "ClientMethod", "ClientProto", "ClientURI",
			"PeerMethod", "PeerProto", "PeerHost",
			"PeerRequestURI", "Description"),
		startTimeLog: time.Now(),
		IsDisabled:   false,
	}

	err := vLog.newFileLog("", true)
	if err != nil {
		return nil, err
	}

	return vLog, nil
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
	go v.addInThread(values...)
}

func (v *VLog) addInThread(values ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch err: %v", err) //nolint:forbidigo
			os.Exit(int(types.ErrGotPanic))
		}
	}()

	if v.IsDisabled || values == nil {
		return
	}

	v.wgNewLog.Wait()

	v.mu.Lock()
	defer v.mu.Unlock()

	if values == nil || v.mapLastLogRecords == nil {
		return
	}

	typeLog, recordRow := v.buildCsvRecord(values)

	//nolint:exhaustive
	switch typeLog {
	case Fatal:
		log.Panic(recordRow)
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

// nolint
func (v *VLog) buildCsvRecord(values []interface{}) (TypeLog, string) {
	var typeLog TypeLog

	var val string

	var resultCode types.ResultCode

	var remoteAddr RemoteAddr

	var clientHost ClientHost

	var clientMethod ClientMethod

	var clientProto ClientProto

	var clientURI ClientURI

	var proxyHost ProxyHost

	var proxyMethod ProxyMethod

	var proxyProto ProxyProto

	var proxyURI ProxyURI

	var isConvertOk bool

	for _, value := range values {
		switch valueTypeLog := value.(type) {
		case TypeLog:
			typeLog, isConvertOk = value.(TypeLog)
			if !isConvertOk {
				continue
			}
		case types.ResultCode:
			resultCode, isConvertOk = value.(types.ResultCode)
			if !isConvertOk {
				continue
			}
		case string:
			val = val + valueTypeLog + ","
		case RemoteAddr:
			remoteAddr, isConvertOk = value.(RemoteAddr)
			if !isConvertOk {
				continue
			}
		case ClientHost:
			clientHost, isConvertOk = value.(ClientHost)
			if !isConvertOk {
				continue
			}
		case ClientMethod:
			clientMethod, isConvertOk = value.(ClientMethod)
			if !isConvertOk {
				continue
			}
		case ClientProto:
			clientProto, isConvertOk = value.(ClientProto)
			if !isConvertOk {
				continue
			}
		case ClientURI:
			clientURI, isConvertOk = value.(ClientURI)
			if !isConvertOk {
				continue
			}
		case ProxyHost:
			proxyHost, isConvertOk = value.(ProxyHost)
			if !isConvertOk {
				continue
			}
		case ProxyMethod:
			proxyMethod, isConvertOk = value.(ProxyMethod)
			if !isConvertOk {
				continue
			}
		case ProxyProto:
			proxyProto, isConvertOk = value.(ProxyProto)
			if !isConvertOk {
				continue
			}
		case ProxyURI:
			proxyURI, isConvertOk = value.(ProxyURI)
			if !isConvertOk {
				continue
			}
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
