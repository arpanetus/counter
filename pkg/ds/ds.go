package ds

import "time"

type RingBufferNoder interface {
	Stamp() time.Time
}

type RingBufferer interface {
	AddNow()
	Sum() (cnt uint64)
	Last() *RingBufferNode
	Add(stamp time.Time)
	Iterator() *RingBufferIterator
}

type RingBufferIterable interface {
	Next() bool
	Value() *RingBufferNode
}
