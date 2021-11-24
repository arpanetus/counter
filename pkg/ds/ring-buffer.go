package ds

import (
	"errors"
	"time"
)

var ErrEmptyRingBuffer = errors.New("ring buffer is empty")

type RingBuffer struct {
	duration time.Duration
	els      map[int64]struct{}
}

// New inits *RingBuffer with given the duration.
func New(duration time.Duration) *RingBuffer {
	return &RingBuffer{duration: duration, els: make(map[int64]struct{})}
}

// AddNow adds current time for a ring.
func (r *RingBuffer) AddNow() {
	r.Add(time.Now())
}

// Sum counts all stamps within RingBuffer's time.Duration,
// and remove the old stamp as soon as it meets it.
func (r *RingBuffer) Sum() (count uint64) {
	now, toRemove := time.Now(), make([]int64, 0)

	for stamp := range r.els {
		if now.Sub(time.Unix(0, stamp)) > r.duration {
			toRemove = append(toRemove, stamp)
		} else {
			count++
		}
	}

	for _, v := range toRemove {
		delete(r.els, v)
	}

	return count
}

// Stamps returns the underlying map.
func (r *RingBuffer) Stamps() map[int64]struct{} {
	return r.els
}

// Add adds the timestamp into *RingBuffer.
func (r *RingBuffer) Add(stamp time.Time) {
	r.els[stamp.UnixNano()] = struct{}{}
}
