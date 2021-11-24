package file

import (
	"io"
	"os"

	"github.com/arpanetus/counter/pkg/ds"
)

type CounterFileWorker interface {
	Read() (b ds.RingBufferer, err error)
	Write(b ds.RingBufferer) error
	Close() error
}

var FS fileSystem = osFS{}

type fileSystem interface {
	Open(name string) (filer, error)
	Create(name string) (filer, error)
	IsNotExist(err error) bool
}

type filer interface {
	io.Closer
	io.Reader
	io.Seeker
	io.StringWriter
	Truncate(size int64) error
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (filer, error)   { return os.Open(name) }
func (osFS) Create(name string) (filer, error) { return os.Create(name) }
func (osFS) IsNotExist(err error) bool         { return os.IsNotExist(err) }
