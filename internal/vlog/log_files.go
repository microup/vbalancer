package vlog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"vbalancer/internal/core"
	"vbalancer/internal/types"
	"vbalancer/internal/version"
)

// newFileLog - function for creating a new file log.
func (v *VLog) newFileLog(isNewFileLog bool) error {
	if isNewFileLog {
		v.MapLastLogRecords = make([]string, 0)

		if err := v.Close(); err != nil {
			return fmt.Errorf("%w", err)
		}
	} else if !isNewFileLog {
		if v.fileLog != nil {
			_ = v.fileLog.Close()
		}
	}

	if _, err := os.Stat(v.cfg.DirLog); os.IsNotExist(err) {
		err = os.Mkdir(v.cfg.DirLog, types.MaskDir)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	var err error
	v.fileLog, err = v.open()
	if v.fileLog == nil || err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err = v.fileLog.WriteString(v.headerCSV + "\n"); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Close function for closing the file log.
func (v *VLog) Close() error {
	v.wg.Wait()

	if v.fileLog != nil {
		if err := v.fileLog.Close(); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

// open function for opening the file log.
func (v *VLog) open() (*os.File, error) {
	var fileNameLog string

	timeCreateLogFile := v.startTimeLog.Format("20060102150405")
	fileNameLog = fmt.Sprintf("%s_%d_%d.%s", timeCreateLogFile, version.Get(), v.idLog, types.LogFileExtension)

	atomic.AddUint64(&v.idLog, 1)

	fileNameLog = filepath.Join(v.cfg.DirLog, fileNameLog)

	fileLog, err := os.OpenFile(fileNameLog, os.O_APPEND|os.O_CREATE|os.O_RDWR, types.DefaultFilePerm)

	if err != nil || fileLog == nil {
		return nil, fmt.Errorf("%w", err)
	}

	return fileLog, nil
}

// GetCurrentFileLogInfo function for getting current file log info.
func (v *VLog) GetCurrentFileLogInfo() *LogFile {
	v.Mu.Lock()
	defer v.Mu.Unlock()

	fileInfo, err := os.Stat(v.fileLog.Name())
	if err != nil {
		return nil
	}

	return &LogFile{
		FileName: v.fileLog.Name(),
		FileSize: core.HumanFileSize(float64(fileInfo.Size())),
	}
}

// checkToCreateNewLogFile function for checking to create a new log file.
func (v *VLog) checkToCreateNewLogFile() error {
	fileInfo, err := os.Stat(v.fileLog.Name())
	if err != nil {
		return fmt.Errorf("failed to os stat: %s err: %w", v.fileLog.Name(), err)
	}

	if fileInfo == nil {
		err = v.newFileLog(false)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	fileSizeBytes := fileInfo.Size()
	fileSizeMB := float64(fileSizeBytes) / (types.LengthKilobytesInBytes * types.LengthKilobytesInBytes)

	if fileSizeMB < v.cfg.FileSizeMB {
		return nil
	}

	oldFileCSV := v.fileLog.Name()

	err = v.newFileLog(false)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = core.ArchiveFile(oldFileCSV, fmt.Sprintf("_%s.zip", types.LogFileExtension))
	if err != nil {
		return fmt.Errorf("failed to archive file: %w", err)
	}

	fileCsv := filepath.Join(v.cfg.DirLog, fileInfo.Name())
	err = os.Remove(fileCsv)

	if err != nil {
		return fmt.Errorf("failed to os remove file:%s err:%w", fileCsv, err)
	}

	return nil
}

// removeOldRecordsFromMemory function for removing old records from memory.
func (v *VLog) removeOldRecordsFromMemory() {
	if uint64(len(v.MapLastLogRecords)) > v.cfg.APIShowRecords {
		var xLast string

		_, v.MapLastLogRecords = v.MapLastLogRecords[0], v.MapLastLogRecords[1:]
		xLast, v.MapLastLogRecords = v.MapLastLogRecords[len(v.MapLastLogRecords)-1],
			v.MapLastLogRecords[:len(v.MapLastLogRecords)-1]

		v.MapLastLogRecords = append(v.MapLastLogRecords, xLast)
	}
}
