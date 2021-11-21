package service

import (
	"fmt"
	"github.com/arpanetus/counter/pkg/ds"
	"github.com/arpanetus/counter/pkg/file"
	"sync"
)

type Counter struct {
	mutex  sync.Mutex
	path   string
	worker file.CounterFileWorker
	buffer ds.RingBufferer
}

func New(
	path string,
	worker file.CounterFileWorker,
) *Counter {
	return &Counter{
		path:   path,
		worker: worker,
	}
}

func (c *Counter) Parse() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	b, err := c.worker.Read(c.path)
	if err != nil {
		return fmt.Errorf("cannot parse from file {%s} with err: %w", c.path, err)
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

	if err := c.worker.Write(c.path, c.buffer); err != nil {
		return fmt.Errorf("cannot write buffer into file {%s} with err", c.path, err)
	}

	return nil
}

func (c *Counter) Close() error {
	return c.Write()
}
