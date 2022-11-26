package vlog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"vbalancer/internal/core"
	"vbalancer/internal/types"
	"vbalancer/internal/version"
)

func (v *VLog) newFileLog(newFileName string, isNewFileLog bool) error {

	if isNewFileLog {
		v.mapLastLogRecords = make([]string, 0)
		err := v.Close()
		if err != nil {
			return err
		}
	} else {
		if v.fileLog != nil {
			_ = v.fileLog.Close()
		}
	}

	var err error
	if _, err = os.Stat(v.cfg.DirLog); os.IsNotExist(err) {
		err = os.Mkdir(v.cfg.DirLog, 0755)
		if err != nil {
			return err
		}
	}

	v.fileLog, err = v.open(newFileName)
	if v.fileLog == nil || err != nil {
		return err
	}

	_, err = v.fileLog.WriteString(v.headerCSV + "\n")

	return err
}

func (v *VLog) Close() error {
	v.wgNewLog.Wait()
	if v.fileLog != nil {
		if err := v.fileLog.Close(); err != nil {
			log.Fatalf("can't close csv log.")
		}
	}
	return nil
}

func (v *VLog) open(newFileName string) (*os.File, error) {

	v.countToLogID = v.countToLogID + 1

	var fileNameLog string
	if newFileName == "" {
		timeCreateLogFile := v.startTimeLog.Format("20060102150405")
		fileNameLog = fmt.Sprintf("%s_%d_%d.%s", timeCreateLogFile, version.Get(), v.countToLogID, v.cfg.KindType)
	} else {
		fileNameLog = fmt.Sprintf("%s_%d_%d.%s", newFileName, version.Get(), v.countToLogID, v.cfg.KindType)
	}

	fileNameLog = filepath.Join(v.cfg.DirLog, fileNameLog)

	fileLog, err := os.OpenFile(fileNameLog, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil || fileLog == nil {
		return nil, fmt.Errorf("CONF: error is opening config file: %s", fileNameLog)
	}

	return fileLog, nil
}

func (v *VLog) GetCurrentFileLogInfo() *types.FileInfo {
	v.mu.Lock()
	defer v.mu.Unlock()

	fileInfo, err := os.Stat(v.fileLog.Name())
	if err != nil {
		return nil
	}

	return &types.FileInfo{
		FileName: v.fileLog.Name(),
		FileSize: core.HumanFileSize(float64(fileInfo.Size())),
		Kind:     v.cfg.KindType,
	}

}

func (v *VLog) checkToCreateNewLogFile() error {

	fileInfo, err := os.Stat(v.fileLog.Name())
	if err != nil {
		return err
	}

	if fileInfo == nil {
		err = v.newFileLog("", false)
		if err != nil {
			return err
		}
		return nil
	}

	if uint64(fileInfo.Size()) > v.cfg.FileSize {

		oldFileCSV := v.fileLog.Name()

		err = v.newFileLog("", false)
		if err != nil {
			return err
		}

		core.ArchiveFile(oldFileCSV, fmt.Sprintf("_%s.zip", v.cfg.KindType))

		fileCsv := filepath.Join(v.cfg.DirLog, fileInfo.Name())
		err = os.Remove(fileCsv)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *VLog) removeOldRecordsFromMemory() {
	if uint64(len(v.mapLastLogRecords)) > v.cfg.ApiShowRecords {
		var xLast string
		_, v.mapLastLogRecords = v.mapLastLogRecords[0], v.mapLastLogRecords[1:]
		xLast, v.mapLastLogRecords = v.mapLastLogRecords[len(v.mapLastLogRecords)-1], v.mapLastLogRecords[:len(v.mapLastLogRecords)-1]
		v.mapLastLogRecords = append(v.mapLastLogRecords, xLast)
	}
}
