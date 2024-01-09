package main

import "time"

const (
	bufSize         int           = 10
	bufEmitInterval time.Duration = 30 * time.Second
)

type data struct {
	value int
	stamp time.Time
}

func NewValue(v int) *data {
	return &data{
		value: v,
		stamp: time.Now(),
	}
}

func (d *data) GetValue() int {
	return d.value
}

func (d *data) GetTimestamp() time.Time {
	return d.stamp
}

type ringBuffer struct {
	data       []*data
	size       int
	lastInsert int
	nextRead   int
	emitTimer  time.Duration
}

func NewRingBuffer(size int, dur time.Duration) *ringBuffer {
	return &ringBuffer{
		data:       make([]*data, size),
		size:       size,
		lastInsert: -1,
		nextRead:   0,
		emitTimer:  dur,
	}
}

func (r *ringBuffer) Insert(input *data) {
	r.lastInsert = (r.lastInsert + 1) % r.size
	r.data[r.lastInsert] = input

	// keep nextRead in front of the lastInsert
	if r.nextRead == r.lastInsert {
		r.nextRead = (r.nextRead + 1) % r.size
	}
}

func (r *ringBuffer) Emit() []*data {
	out := []*data{}
	for {
		if r.data[r.nextRead] != nil {
			out = append(out, r.data[r.nextRead])
			r.data[r.nextRead] = nil // clean the data structure once we have read it
		}
		if r.nextRead == r.lastInsert || r.lastInsert == -1 {
			break
		}
		r.nextRead = (r.nextRead + 1) % r.size
	}
	return out
}
