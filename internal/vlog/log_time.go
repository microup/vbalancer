package vlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"vbalancer/internal/core"
	"vbalancer/internal/types"
)

func (v *VLog) New(fileNameDateTime string) error {
	defer func() {
		v.startTimeLog = time.Now()
	}()

	v.wgNewLog.Add(1)

	defer func() {
		v.wgNewLog.Done()
	}()

	var fileInfo os.FileInfo

	var err error

	if v.fileLog != nil {
		fileInfo, err = os.Stat(v.fileLog.Name())
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	if fileInfo == nil {
		err = v.newFileLog(fileNameDateTime, false)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		return types.ErrFileIsNil
	}

	oldFileCSV := v.fileLog.Name()

	err = v.newFileLog(fileNameDateTime, false)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = core.ArchiveFile(oldFileCSV, fmt.Sprintf("_%s.zip", types.LogFileExtension))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fileCsv := filepath.Join(v.cfg.DirLog, fileInfo.Name())
	err = os.Remove(fileCsv)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
