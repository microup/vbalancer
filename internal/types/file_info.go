package types

type FileInfo struct {
	FileName string `json:"fileName,omitempty"`
	FileSize string `json:"fileSize,omitempty"`
	Kind     string `json:"kind,omitempty"`
}

type UpdateFileInfo func() *FileInfo
