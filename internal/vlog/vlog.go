package vlog

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"vbalancer/internal/config"
	"vbalancer/internal/types"
)

type VLog struct {
	cfg               *config.Log
	fileLog           *os.File
	countToLogID      int
	MapLastLogRecords []string
	Mu                *sync.Mutex
	headerCSV         string
	startTimeLog      time.Time
	wgNewLog          *sync.WaitGroup
	IsDisabled        bool
}

func New(cfg *config.Log) *VLog {
	headerCSV := fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;",
		"Date",
		"Time",
		"Type",
		"ResultCode",
		"RemoteAddr",
		"PeerURL",
		"Description")

	return &VLog{
		wgNewLog:          &sync.WaitGroup{},
		Mu:                &sync.Mutex{},
		fileLog:           nil,
		cfg:               cfg,
		countToLogID:      -1,
		MapLastLogRecords: []string{},
		headerCSV:         headerCSV,
		startTimeLog:      time.Now(),
		IsDisabled:        false,
	}
}

func (v *VLog) Init() error {
	err := v.newFileLog("", true)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
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
			log.Fatalf("error recovery: %d, catch err: %s", types.ErrRecoverPanic,  err)
		}
	}()

	if v.IsDisabled || values == nil {
		return
	}

	v.wgNewLog.Wait()

	v.Mu.Lock()
	defer v.Mu.Unlock()

	if v.MapLastLogRecords == nil {
		return
	}

	typeLog, recordRow := BuildRecord(types.ParseValues(values))

	if typeLog == types.Fatal {
		log.Panic(recordRow)
	}

	log.Print(recordRow)

	_, err := v.fileLog.WriteString(recordRow + "\n")
	if err != nil {
		log.Printf("Error: %s - is writing: %s to log file: %s\n", err, recordRow, v.fileLog.Name())

		return
	}

	v.removeOldRecordsFromMemory()

	v.MapLastLogRecords = append(v.MapLastLogRecords, recordRow)

	err = v.checkToCreateNewLogFile()
	if err != nil {
		log.Printf("Error: %s - is writing: %s to log file: %s", err, recordRow, v.fileLog.Name())
	}
}
