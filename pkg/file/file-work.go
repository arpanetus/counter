package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/arpanetus/counter/pkg/ds"
)

type CounterFileWork struct {
	file filer
	fs   fileSystem
	path string
	dur  time.Duration
}

func New(path string, dur time.Duration, fs fileSystem) *CounterFileWork {
	return &CounterFileWork{path: path, dur: dur, fs: fs}
}

func (w *CounterFileWork) parser(b ds.RingBufferer, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		i, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse line: %w", err)
		}

		b.Add(time.Unix(0, i))
	}
	return nil
}

func (w *CounterFileWork) Read() (ds.RingBufferer, error) {
	f, err := w.fs.Open(w.path)
	if err != nil && !w.fs.IsNotExist(err) {
		return nil, fmt.Errorf("cannot open f: %w", err)
	}
	if os.IsNotExist(err) {
		f, err = w.fs.Create(w.path)
		if err != nil {
			return nil, fmt.Errorf("cannot create file: %w", err)
		}
	}

	b := ds.New(w.dur)
	if err := w.parser(b, f); err != nil {
		return nil, fmt.Errorf("cannot parse file: %w", err)
	}

	if err = f.Close(); err != nil {
		return nil, fmt.Errorf("cannot close file: %w", err)
	}

	return b, nil
}

func (w *CounterFileWork) Write(b ds.RingBufferer) error {
	if b == nil {
		return ds.ErrEmptyRingBuffer
	}
	if stamps := b.Stamps(); len(stamps) == 0 {
		return ds.ErrEmptyRingBuffer
	}

	if w.file == nil {
		f, err := w.fs.Create(w.path)
		if err != nil {
			return fmt.Errorf("cannot open file for creation: %w", err)
		}
		w.file = f
	}

	if err := w.file.Truncate(0); err != nil {
		fmt.Printf("cannot truncate: %+v", err)
	}
	if _, err := w.file.Seek(0, 0); err != nil {
		fmt.Printf("cannot seek: %+v", err)
	}

	for stamp := range b.Stamps() {
		line := strconv.FormatInt(stamp, 10) + "\n"

		if _, err := w.file.WriteString(line); err != nil {
			return fmt.Errorf("cannot append line {%s} into file: %w", line, err)
		}
	}

	return nil
}

func (w *CounterFileWork) Close() error {
	if err := w.file.Close(); err != nil {
		return fmt.Errorf("cannot close file")
	}
	return nil
}
