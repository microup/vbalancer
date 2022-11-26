package vlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"vbalancer/internal/core"
)

func (v *VLog) New(fileNameDateTime string) {

	defer func() {
		v.startTimeLog = time.Now()
	}()

	v.wgNewLog.Add(1)
	defer func() {
		v.wgNewLog.Done()
	}()

	var fileInfo os.FileInfo = nil
	var err error
	if v.fileLog != nil {
		fileInfo, err = os.Stat(v.fileLog.Name())
		if err != nil {
			return
		}
	}

	if fileInfo == nil {
		err = v.newFileLog(fileNameDateTime, false)
		if err != nil {
			return
		}
		return
	}

	oldFileCSV := v.fileLog.Name()

	err = v.newFileLog(fileNameDateTime, false)
	if err != nil {
		return
	}

	core.ArchiveFile(oldFileCSV, fmt.Sprintf("_%s.zip", v.cfg.KindType))

	fileCsv := filepath.Join(v.cfg.DirLog, fileInfo.Name())
	err = os.Remove(fileCsv)
	if err != nil {
		return
	}
}
