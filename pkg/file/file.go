package file

import "github.com/arpanetus/counter/pkg/ds"

type CounterFileWorker interface {
	Read(path string) (b ds.RingBufferer, err error)
	Write(path string, b ds.RingBufferer) error
}
