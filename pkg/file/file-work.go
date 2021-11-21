package file

import (
	"bufio"
	"fmt"
	"github.com/arpanetus/counter/pkg/ds"
	"log"
	"os"
	"time"
)

type CounterFileWork struct {
	dur        time.Duration
	timeFormat string
}

func New(dur time.Duration, timeFormat string) *CounterFileWork {
	return &CounterFileWork{dur: dur, timeFormat: timeFormat}
}

func (w CounterFileWork) Read(path string) (b ds.RingBufferer, err error) {


	f, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot open f: %w", err)
	}
	if os.IsNotExist(err) {
		f, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("cannot create file: %w", err)
		}
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("cannot close f: %v", err)
		}
	}(f)

	b = ds.New(w.dur)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		t, err := time.Parse(w.timeFormat, line)
		if err != nil {
			return nil, fmt.Errorf("cannot parse time: %w", err)
		}
		b.Add(t)
	}

	return b, nil
}

func (w CounterFileWork) Write(path string, b ds.RingBufferer) error {
	if b == nil {
		return ds.ErrEmptyRingBuffer
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("cannot close file: %v", err)
		}
	}(f)

	i := b.Iterator()

	for i.Next() {
		line := i.Value().Stamp().Format(w.timeFormat) + "\n"

		_, err = f.WriteString(line)
		if err != nil {
			return fmt.Errorf("cannot append line {%s} into file: %w", line, err)
		}
	}

	return nil
}
