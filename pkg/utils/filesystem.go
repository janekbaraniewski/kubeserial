package utils

import (
	"io"
	"os"

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
type inMemoryFS struct {
	fs afero.Fs
}

func (f *inMemoryFS) Open(name string) (File, error)        { return f.fs.Open(name) }
func (f *inMemoryFS) Stat(name string) (os.FileInfo, error) { return f.fs.Stat(name) }
func (f *inMemoryFS) Create(name string) (File, error)      { return f.fs.Create(name) }

func NewInMemoryFS() *inMemoryFS {
	fs := &inMemoryFS{}
	fs.fs = afero.NewMemMapFs()
	return fs
}
