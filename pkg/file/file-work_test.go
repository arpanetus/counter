package file

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/arpanetus/counter/pkg/ds"
)

func assertEqual(t *testing.T, got interface{}, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s != %s", got, want)
	}
}

const (
	emptyPath = ""
	dur       = time.Minute
)

type MockRing struct {
	els []time.Time
}

func (m *MockRing) Add(t time.Time) {
	m.els = append(m.els, t)
}
func (*MockRing) AddNow()                    {}
func (*MockRing) Sum() uint64                { return 0 }
func (*MockRing) Stamps() map[int64]struct{} { return nil }

var (
	ReaderMockOk    = strings.NewReader("1637698562645791724\n1637700320092426481\n")
	ReaderMockNonOk = strings.NewReader("HereCouldBeYourAdd:)\n")
	Times           = []time.Time{time.Unix(0, 1637698562645791724), time.Unix(0, 1637700320092426481)}
	ErrSome         = errors.New("some error")
)

func TestCounterFileWork_parser(t *testing.T) {
	fw, r := New(emptyPath, dur, FS), &MockRing{els: make([]time.Time, 0)}

	assertEqual(t, fw.parser(r, ReaderMockOk), nil)
	assertEqual(t, r.els, Times)

	fw, r = New(emptyPath, dur, FS), &MockRing{els: make([]time.Time, 0)}

	err := fw.parser(r, ReaderMockNonOk)
	if err == nil {
		t.Fatalf("%s != nil", err)
	}
	assertEqual(t, r.els, make([]time.Time, 0))
	var match *strconv.NumError
	if !errors.As(err, &match) {
		t.Fatalf("%s != %s", match, err)
	}
}

type file struct {
	cnt     int
	written string
}

func (f *file) Read(p []byte) (n int, err error) {
	if f.cnt == 1 {
		return 0, ErrSome
	}
	return ReaderMockOk.Read(p)
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *file) Close() error {
	if f.cnt == 1 {
		return ErrSome
	}

	return nil
}

func (f *file) WriteString(s string) (n int, err error) {
	if f.cnt == 1 {
		return 0, ErrSome
	}
	f.written = s
	return len([]byte(s)), nil
}

func (f *file) Truncate(size int64) error {
	return nil
}

type fs struct {
	cnt int
	filer
}

func (f *fs) Open(name string) (filer, error) {
	if f.cnt == 1 {
		return nil, os.ErrNotExist
	}
	return f.filer, nil
}

func (f *fs) Create(name string) (filer, error) {
	if f.cnt == 1 {
		return nil, &os.PathError{Op: "open", Path: name, Err: ErrSome}
	}

	return f.filer, nil
}

func (f *fs) IsNotExist(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

func TestCounterFileWork_ReadWriteClose(t *testing.T) {
	var mockFile file
	mockFs := fs{cnt: 0, filer: &mockFile}
	fw := New(emptyPath, dur, &mockFs)

	b, err := fw.Read()
	if errors.Is(err, os.ErrClosed) {
		t.Fatalf("%s!=%s", err, os.ErrClosed)
	}

	err = fw.Write(b)
	if errors.Is(err, os.ErrClosed) {
		t.Fatalf("%s!=%s", err, os.ErrClosed)
	}

	err = fw.Close()
	if errors.Is(err, os.ErrClosed) {
		t.Fatalf("%s!=%s", err, os.ErrClosed)
	}
}

func TestCounterFileWork_ReadWriteCloseFail(t *testing.T) {
	var mockFile file
	mockFile.cnt = 1
	mockFs := fs{cnt: 1, filer: &mockFile}
	fw := New(emptyPath, dur, &mockFs)

	b, err := fw.Read()
	assertEqual(t, b, nil)
	if !errors.Is(err, ErrSome) {
		t.Fatalf("%s!=%s", err, ErrSome)
	}

	err = fw.Write(nil)
	if !errors.Is(err, ds.ErrEmptyRingBuffer) {
		t.Fatalf("%s!=%s", err, ds.ErrEmptyRingBuffer)
	}

	err = fw.Write(&ds.RingBuffer{})
	if !errors.Is(err, ds.ErrEmptyRingBuffer) {
		t.Fatalf("%s!=%s", err, ds.ErrEmptyRingBuffer)
	}
}
