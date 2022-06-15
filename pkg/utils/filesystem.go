package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// FileSystem interface describes FS methods used inside KubeSerial implementation.
type FileSystem interface {
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
}

// File interface describes File structure.
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	Stat() (os.FileInfo, error)
}

// OsFS implements fileSystem using the local disk.
type OsFS struct{}

// Open takes file name and returns File from filesystem.
func (f *OsFS) Open(name string) (File, error) { return os.Open(name) }

// Stat takes file name and returns os.FileInfo of File.
func (f *OsFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

// NewOSFS returns new instance of OsFS.
func NewOSFS() *OsFS { return &OsFS{} }

// InMemoryFS implements fileSystem using afero.
type InMemoryFS struct {
	fs afero.Fs
}

// Open takes file name and returns File from filesystem.
func (f *InMemoryFS) Open(name string) (File, error) { return f.fs.Open(name) }

// Stat takes file name and returns os.FileInfo of File.
func (f *InMemoryFS) Stat(name string) (os.FileInfo, error) { return f.fs.Stat(name) }

// Create takes file name, creates File under this name and returns it for usage.
func (f *InMemoryFS) Create(name string) (File, error) { return f.fs.Create(name) }

// AddFileFromHostPath takes file path of file inside test-assets and places it inside /config dir in mem fs.
func (f *InMemoryFS) AddFileFromHostPath(path string) error {
	_, fileName := filepath.Split(path)
	file, err := f.Create(path)
	if err != nil {
		return err
	}

	absPath, _ := filepath.Abs(fmt.Sprintf("../../test-assets/%v", fileName))
	content, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	file.Close()
	return nil
}

// NewInMemoryFS returns new instance of InMemoryFS.
func NewInMemoryFS() *InMemoryFS {
	fs := &InMemoryFS{}
	fs.fs = afero.NewMemMapFs()
	return fs
}
