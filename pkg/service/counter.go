package service

import (
	"fmt"
	"sync"

	"github.com/arpanetus/counter/pkg/ds"
	"github.com/arpanetus/counter/pkg/file"
)

type Counter struct {
	mutex  sync.Mutex
	worker file.CounterFileWorker
	buffer ds.RingBufferer
}

func New(
	worker file.CounterFileWorker,
) *Counter {
	return &Counter{
		worker: worker,
	}
}

func (c *Counter) Parse() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	b, err := c.worker.Read()
	if err != nil {
		return fmt.Errorf("cannot parse from file: %w", err)
	}

	c.buffer = b

	return nil
}

func (c *Counter) Count() uint64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.buffer.Sum()
}

func (c *Counter) Add() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.buffer.AddNow()
}

func (c *Counter) Write() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.worker.Write(c.buffer); err != nil {
		return fmt.Errorf("cannot write buffer into file: %w", err)
	}

	return nil
}

func (c *Counter) Close() error {
	return c.worker.Close()
}
