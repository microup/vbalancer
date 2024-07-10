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
	wg                *sync.WaitGroup
	cfg               *config.Log
	fileLog           *os.File
	idLog             uint64
	MapLastLogRecords []string
	Mu                *sync.Mutex
	headerCSV         string
	startTimeLog      time.Time
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
		Mu:                &sync.Mutex{},
		wg:                &sync.WaitGroup{},
		fileLog:           nil,
		cfg:               cfg,
		idLog:             0,
		MapLastLogRecords: []string{},
		headerCSV:         headerCSV,
		startTimeLog:      time.Now(),
		IsDisabled:        false,
	}
}

func (v *VLog) Init() error {
	if err := v.newFileLog(true); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// GetCountRecords returns the number of records in the log file.
func (v *VLog) GetCountRecords() int {
	v.Mu.Lock()
	defer v.Mu.Unlock()

	return len(v.MapLastLogRecords)
}

// Add adds a log record to the log file.
func (v *VLog) Add(values ...interface{}) {
	v.wg.Add(1)

	go v.addInThread(values...)
}

func (v *VLog) addInThread(values ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("error recovery: %d, catch err: %s", types.ErrRecoverPanic, err)
		}
	}()

	if v.IsDisabled || values == nil {
		return
	}

	v.Mu.Lock()
	defer v.Mu.Unlock()

	defer v.wg.Done()

	typeLog, recordRow := BuildRecord(types.ParseValues(values))

	if typeLog == types.Fatal {
		log.Panic(recordRow)
	}

	log.Print(recordRow)

	if _, err := v.fileLog.WriteString(recordRow + "\n"); err != nil {
		log.Printf("Error: %s - is writing: %s to log file: %s\n", err, recordRow, v.fileLog.Name())

		return
	}

	v.removeOldRecordsFromMemory()

	v.MapLastLogRecords = append(v.MapLastLogRecords, recordRow)

	if err := v.checkToCreateNewLogFile(); err != nil {
		log.Printf("Error: %s - is writing: %s to log file: %s", err, recordRow, v.fileLog.Name())
	}
}
