package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type FileSystem interface {
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	Stat() (os.FileInfo, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (f *osFS) Open(name string) (File, error)        { return os.Open(name) }
func (f *osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

func NewOSFS() *osFS { return &osFS{} }

// aferoFS implements fileSystem using afero
type InMemoryFS struct {
	fs afero.Fs
}

func (f *InMemoryFS) Open(name string) (File, error)        { return f.fs.Open(name) }
func (f *InMemoryFS) Stat(name string) (os.FileInfo, error) { return f.fs.Stat(name) }
func (f *InMemoryFS) Create(name string) (File, error)      { return f.fs.Create(name) }
func (f *InMemoryFS) AddFileFromHostPath(path string) error {
	file, err := f.Create(fmt.Sprintf("/config/%v", path))

	if err != nil {
		return err
	}

	absPath, _ := filepath.Abs(fmt.Sprintf("../../test-assets/%v", path))
	content, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	file.Write(content)
	file.Close()
	return nil
}

func NewInMemoryFS() *InMemoryFS {
	fs := &InMemoryFS{}
	fs.fs = afero.NewMemMapFs()
	return fs
}
