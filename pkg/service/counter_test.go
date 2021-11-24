package service

import (
	"github.com/arpanetus/counter/pkg/ds"
	"reflect"
	"testing"
	"time"
)

func assertEqual(t *testing.T, got interface{}, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s != %s", got, want)
	}
}

const (
	dur = time.Minute
)

var (
	time1 = time.Unix(0, 1637698562645791724)
	time2 = time.Unix(0, 1637700320092426481)
)

type FileWorkerMock struct {
	written ds.RingBufferer
}

func (f *FileWorkerMock) Read() (b ds.RingBufferer, err error) {
	b = ds.New(dur)
	b.Add(time1)
	b.Add(time2)
	return b, nil
}

func (f *FileWorkerMock) Write(b ds.RingBufferer) error {
	f.written = b
	return nil
}

func (f *FileWorkerMock) Close() error {
	return nil
}

func TestCounter_ParseAdd(t *testing.T) {
	b := ds.New(dur)
	b.Add(time1)
	b.Add(time2)

	fw := FileWorkerMock{}
	svc := New(&fw)
	err := svc.Parse()
	if err != nil {
		t.Fatalf("%s!=nil", err)
	}
	assertEqual(t, svc.buffer, b)

	svc.Add()
	assertEqual(t, len(svc.buffer.Stamps()), 3)

	res := svc.Count()
	assertEqual(t, res, uint64(1))

	err = svc.Write()
	if err != nil {
		t.Fatalf("%s!=nil", err)
	}

	assertEqual(t, len(fw.written.Stamps()), 1)
}
