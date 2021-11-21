package ds

import (
	"errors"
	"time"
)

var ErrEmptyRingBuffer = errors.New("ring buffer is empty")

type RingBufferNode struct {
	stamp time.Time
	next  *RingBufferNode
}

func (n *RingBufferNode) Stamp() time.Time {
	return n.stamp
}

type RingBuffer struct {
	duration  time.Duration
	startNode *RingBufferNode
	lastNode  *RingBufferNode
}

func emptyRingBufferNode() *RingBufferNode {
	return &RingBufferNode{stamp: time.Time{}, next: nil}
}

// New inits *RingBuffer with given the duration and zero-value *RingBuffer as the startNode and lastNode.
func New(duration time.Duration) *RingBuffer {
	node := emptyRingBufferNode()
	return &RingBuffer{duration: duration, startNode: node, lastNode: node}
}

// AddNow adds current time for a ring.
func (r *RingBuffer) AddNow() {
	r.Add(time.Now())
}

// Sum counts all stamps within RingBuffer's time.Duration,
// and updates as soon as it meets the old stamp.
func (r *RingBuffer) Sum() (cnt uint64) {
	node := r.startNode
	moved := false

	for node != nil {
		if node.stamp.Sub(time.Now()) <= r.duration {
			if !moved {
				r.startNode = node
				moved = true
			}
			cnt += 1
		}
		node = node.next
	}

	return cnt
}

func (r *RingBuffer) Last() *RingBufferNode {
	switch {
	case r.lastNode != nil:
		return r.lastNode
	case r.startNode != nil && r.lastNode == nil:
		node := r.startNode
		for node.next != nil {
			return node
		}
	default:
		return nil
	}
	return nil
}

// Add add the timestamp into *RingBuffer.
func (r *RingBuffer) Add(stamp time.Time) {
	now := &RingBufferNode{stamp: stamp, next: nil}
	switch {
	case r.startNode == nil && r.lastNode == nil:
		r.startNode, r.lastNode = now, now
	case r.startNode == nil && r.lastNode != nil:
		r.startNode = r.lastNode
		r.lastNode.next = now
		r.lastNode = r.lastNode.next
	case r.startNode != nil && r.lastNode == nil:
		lastNode := r.Last()
		if lastNode == nil {
			r.startNode, r.lastNode = now, now
		} else {
			r.lastNode = lastNode
			r.lastNode.next = now
			r.lastNode = r.lastNode.next
		}
	default:
		r.lastNode.next = now
		r.lastNode = r.lastNode.next
	}
}

type RingBufferIterator struct {
	node  *RingBufferNode
	moved bool
}

func (r *RingBuffer) Iterator() *RingBufferIterator {
	return &RingBufferIterator{
		node:  r.startNode,
		moved: false,
	}
}

func (i *RingBufferIterator) Next() bool {
	if !i.moved {
		i.moved = true

		return true
	}

	if i.node == nil {
		return false
	}

	return true
}

func (i *RingBufferIterator) Value() *RingBufferNode {
	return i.node
}
