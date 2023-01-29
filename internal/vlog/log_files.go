package vlog

import (
	"fmt"
	"os"
	"path/filepath"
	"vbalancer/internal/core"
	"vbalancer/internal/types"
	"vbalancer/internal/version"
)

const (
	maskDir                     = 0x755
	DefaultFilePerm os.FileMode = 0666
)

func (v *VLog) newFileLog(newFileName string, isNewFileLog bool) error {
	if isNewFileLog {
		v.mapLastLogRecords = make([]string, 0)
		err := v.Close()

		if err != nil {
			return err
		}
	} else if !isNewFileLog {
		if v.fileLog != nil {
			_ = v.fileLog.Close()
		}
	}

	var err error
	if _, err = os.Stat(v.cfg.DirLog); os.IsNotExist(err) {
		err = os.Mkdir(v.cfg.DirLog, maskDir)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	v.fileLog, err = v.open(newFileName)
	if v.fileLog == nil || err != nil {
		return err
	}

	_, err = v.fileLog.WriteString(v.headerCSV + "\n")
	if err != nil {
		return fmt.Errorf("%w", err)
	} else {
		return nil
	}

}

func (v *VLog) Close() error {
	v.wgNewLog.Wait()

	if v.fileLog != nil {
		if err := v.fileLog.Close(); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func (v *VLog) open(newFileName string) (*os.File, error) {
	v.countToLogID++

	var fileNameLog string

	if newFileName == "" {
		timeCreateLogFile := v.startTimeLog.Format("20060102150405")
		fileNameLog = fmt.Sprintf("%s_%d_%d.%s", timeCreateLogFile, version.Get(), v.countToLogID, types.LogFileExtension)
	} else {
		fileNameLog = fmt.Sprintf("%s_%d_%d.%s", newFileName, version.Get(), v.countToLogID, types.LogFileExtension)
	}

	fileNameLog = filepath.Join(v.cfg.DirLog, fileNameLog)

	//nolint:nosnakecase
	fileLog, err := os.OpenFile(fileNameLog, os.O_APPEND|os.O_CREATE|os.O_RDWR, DefaultFilePerm)

	if err != nil || fileLog == nil {
		return nil, fmt.Errorf("%w", err)
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
		Kind:     types.LogFileExtension,
	}
}

func (v *VLog) checkToCreateNewLogFile() error {
	fileInfo, err := os.Stat(v.fileLog.Name())
	if err != nil {
		return fmt.Errorf("failed to os stat: %s err: %w", v.fileLog.Name(), err)
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

		core.ArchiveFile(oldFileCSV, fmt.Sprintf("_%s.zip", types.LogFileExtension))

		fileCsv := filepath.Join(v.cfg.DirLog, fileInfo.Name())
		err = os.Remove(fileCsv)

		if err != nil {
			return fmt.Errorf("failed to os remove file:%s err:%w", fileCsv, err)
		}
	}

	return nil
}

func (v *VLog) removeOldRecordsFromMemory() {
	if uint64(len(v.mapLastLogRecords)) > v.cfg.APIShowRecords {
		var xLast string

		_, v.mapLastLogRecords = v.mapLastLogRecords[0], v.mapLastLogRecords[1:]
		xLast, v.mapLastLogRecords = v.mapLastLogRecords[len(v.mapLastLogRecords)-1],
			v.mapLastLogRecords[:len(v.mapLastLogRecords)-1]
		v.mapLastLogRecords = append(v.mapLastLogRecords, xLast)
	}
}
