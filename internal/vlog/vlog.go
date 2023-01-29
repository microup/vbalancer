package vlog

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"vbalancer/internal/types"
)

type ILog interface {
	Add(values ...interface{})
	Close() error
}

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

	typeLog, recordRow := types.BuildRecord(types.ParseValues(values))

	//nolint:exhaustive
	switch typeLog {
	case types.Fatal:
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
