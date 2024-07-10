package vlog

// LogFile represents information about a file, including its name, size, and kind.
type LogFile struct {
	// FileName is the name of the file.
	FileName string `json:"fileName,omitempty"`
	// FileSize is the size of the file in string format.
	FileSize string `json:"fileSize,omitempty"`
}
