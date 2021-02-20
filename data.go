package requests

import (
	"path/filepath"
)

type (
	Auth map[string]interface{}

	Form map[string]interface{}

	Payload map[string]interface{}

	File struct {
		Src      []byte
		Param    string
		Path     string
		Name     string
		MimeType string
		Args     map[string]string
	}
)

func (d *Auth) Update(s Auth) {
	for k, v := range s {
		(*d)[k] = v
	}
}

func (d *Form) Update(s Form) {
	for k, v := range s {
		(*d)[k] = v
	}
}

func (d *Payload) Update(s Payload) {
	for k, v := range s {
		(*d)[k] = v
	}
}

func FileFromBytes(filename string, src []byte) *File {
	return &File{
		Src:  src,
		Name: filename,
	}
}

func FileFromPath(path string) *File {
	return &File{
		Path: path,
		Name: filepath.Base(path),
	}
}

func (f *File) SetName(filename string) *File {
	f.Name = filename
	return f
}

func (f *File) MIME(mimeType string) *File {
	f.MimeType = mimeType
	return f
}
