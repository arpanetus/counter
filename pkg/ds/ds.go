package ds

import "time"

// RingBufferer is the interface which is made for file worker.
// In fact, the prior implementation was a linked-list as a Ring Buffer.
// The better implementation would be a heap.
type RingBufferer interface {
	AddNow()
	Sum() (count uint64)
	Add(stamp time.Time)
	Stamps() map[int64]struct{}
}
