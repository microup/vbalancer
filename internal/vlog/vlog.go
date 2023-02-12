package vlog

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"vbalancer/internal/types"
)

// Ilog is the interface for log. 
type ILog interface {
	Add(values ...interface{})
	Close() error
}

type VLog struct {
	cfg               *Config
	fileLog           *os.File
	countToLogID      int
	MapLastLogRecords []string
	Mu                *sync.Mutex
	headerCSV         string
	startTimeLog      time.Time
	wgNewLog          *sync.WaitGroup
	IsDisabled        bool
}

func New(cfg *Config) (*VLog, error) {
	vLog := &VLog{
		wgNewLog:          &sync.WaitGroup{},
		Mu:                &sync.Mutex{},
		fileLog:           nil,
		cfg:               cfg,
		countToLogID:      -1,
		MapLastLogRecords: []string{},
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

// GetCountRecords returns the number of records in the log file.
func (v *VLog) GetCountRecords() int {
	v.Mu.Lock()
	defer v.Mu.Unlock()

	if v.MapLastLogRecords == nil {
		return 0
	}

	return len(v.MapLastLogRecords)
}

// Add adds a log record to the log file.
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

	v.Mu.Lock()
	defer v.Mu.Unlock()

	if values == nil || v.MapLastLogRecords == nil {
		return
	}

	typeLog, recordRow := BuildRecord(ParseValues(values))

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

	v.MapLastLogRecords = append(v.MapLastLogRecords, recordRow)

	err = v.checkToCreateNewLogFile()
	if err != nil {
		log.Printf("Error: %s - is writing: %s to log file: %s", err, recordRow, v.fileLog.Name())
	}
}
