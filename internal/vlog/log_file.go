package vlog

// LogFile represents information about a file, including its name, size, and kind.
type LogFile struct {
	// FileName is the name of the file.
	FileName string `json:"fileName,omitempty"`
	// FileSize is the size of the file in string format.
	FileSize string `json:"fileSize,omitempty"`
	// Kind is the type of file.
	Kind string `json:"kind,omitempty"`
}

// UpdateFileInfo is a type that represents a function that returns a pointer to a FileInfo.
type UpdateFileInfo func() *LogFile
