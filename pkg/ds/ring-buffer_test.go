package ds

import (
	"reflect"
	"testing"
	"time"
)

func assertEqual(t *testing.T, got interface{}, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s != %s", got, want)
	}
}

const d = time.Minute

func TestNew(t *testing.T) {
	assertEqual(t, New(d), &RingBuffer{duration: d, els: make(map[int64]struct{})})
}

func TestRingBuffer_Add(t *testing.T) {
	b := New(d)
	a := time.Now()
	b.Add(a)
	for i := range b.Stamps() {
		assertEqual(t, a.UnixNano(), i)
	}
}

func TestRingBuffer_AddNow(t *testing.T) {
	b := New(d)
	b.AddNow()
	for i := range b.Stamps() {
		// Of course, it's a test for the sake of testing :).
		want := time.Now()
		got := time.Unix(0, i)
		if want.Sub(got) > time.Millisecond {
			t.Fatalf("your machine seems to be veeeery slow since %s is less than %s", want, got)
		}
	}
}

func TestRingBuffer_Stamps(t *testing.T) {
	b := New(d)
	a := time.Now()
	b.Add(a)
	assertEqual(t, b.Stamps(), map[int64]struct{}{a.UnixNano(): {}})
}

func TestRingBuffer_Sum(t *testing.T) {
	//  Kinda doing stress test, it can't go beyond 1b.
	b := New(d)
	num := 1000000
	for i := 0; i < num; i++ {
		b.AddNow()
	}
	assertEqual(t, b.Sum(), uint64(num))

	// Check for moving window.
	b = New(d)
	num = 10000
	for i := 0; i < num/2; i++ {
		b.Add(time.Now().Add(-d))
	}
	for i := 0; i < num/2; i++ {
		b.AddNow()
	}
	assertEqual(t, b.Sum(), uint64(5000))

	assertEqual(t, len(b.Stamps()), 5000)
}
