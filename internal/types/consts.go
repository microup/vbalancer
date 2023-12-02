package types

import (
	"os"
	"time"
)

// CORE - CONSTS.
const (
	// LengthByte is the length of a byte.
	LengthByte = 1024
	//
	LengthKilobytesInBytes = 1024
	// PowX is the power of 10.
	PowX = 10
	// RoundOne is the value of 0.5.
	RoundOne = .5
	// BitSize is the size of a bit.
	BitSize = 64
)

// CONFIG - CONSTS.
const (
	//
	DefaultNameConfigFile = "config.yaml"
	// MaskDir 0x755 is an octal notation for the file permission -rwxr-xr-x.
	MaskDir = 0x755
	// DefaultFilePerm is the default file permission with octal notation 0666.
	DefaultFilePerm os.FileMode = 0666
	// DefaultFileHeader is the header row for log files.
	DefaultFileLogSizeMB = 10
	// DeafultShowRecordsAPI is the default number of records to show via API.
	DeafultShowRecordsAPI = 50
	// headerCSV defines the header of the CSV log file.
	DefaultDirLogs = "/logs"
	// LogFileExtension is the file extension used for log files in CSV format.
	LogFileExtension = "csv"
)

// PROXY - CONSTS.
const (
	//
	DefaultProxyPort = "8080"
	//
	DeafultMaxCountConnection = 1000
	//
	DeafultClientDeadLineTime = 30 * time.Second
	//
	DeafultPeerConnectionTimeout = 30 * time.Second
	//
	DeafultPeerHostDeadLine = 30 * time.Second
	//
	DeafultCountMaxDialAttemptsToPeer = 30
)
