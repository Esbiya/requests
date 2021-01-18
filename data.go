package requests

import (
	"path/filepath"
)

type (
	Auth map[string]interface{}

	Form map[string]interface{}

	Payload map[string]interface{}

	FileOption struct {
		Src       []byte
		FileParam string
		FilePath  string
		FileName  string
		MimeType  string
	}

	Files map[string]interface{}
)

func File(filename string, src []byte) *FileOption {
	return &FileOption{
		Src:      src,
		FileName: filename,
	}
}

func FileFromPath(path string) *FileOption {
	return &FileOption{
		FilePath: path,
		FileName: filepath.Base(path),
	}
}

func (f *FileOption) FName(filename string) *FileOption {
	f.FileName = filename
	return f
}

func (f *FileOption) MIME(mimeType string) *FileOption {
	f.MimeType = mimeType
	return f
}
