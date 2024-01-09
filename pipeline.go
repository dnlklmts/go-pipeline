package main

import (
	"time"
)

type stageInts func(done <-chan bool, data <-chan int) <-chan int

type pipeline struct {
	stages []stageInts
	done   <-chan bool
}

func NewPipeline(done <-chan bool, stages ...stageInts) *pipeline {
	return &pipeline{
		stages: stages,
		done:   done,
	}
}

func (p *pipeline) Run(src <-chan int) <-chan int {
	var ch <-chan int = src
	for _, stage := range p.stages {
		ch = p.runStage(stage, ch)
	}
	return ch
}

func (p *pipeline) runStage(stage stageInts, srcChan <-chan int) <-chan int {
	return stage(p.done, srcChan)
}

func filterNegativeValues(done <-chan bool, data <-chan int) <-chan int {
	filtered := make(chan int)

	go func() {
		for val := range data {
			// skip negative values
			if val < 0 {
				continue
			}

			select {
			case <-done:
				return
			case filtered <- val:
			}
		}
	}()

	return filtered
}

func filterSpecificValues(done <-chan bool, data <-chan int) <-chan int {
	filtered := make(chan int)

	go func() {
		for val := range data {
			// skip 0 and not divisible by three
			if val == 0 || (val%3 != 0) {
				continue
			}

			select {
			case <-done:
				return
			case filtered <- val:
			}
		}
	}()

	return filtered
}

func bufferValues(done <-chan bool, data <-chan int) <-chan int {
	buffered := make(chan int)
	b := NewRingBuffer(bufSize, bufEmitInterval)

	go func() {
		for {
			select {
			case <-done:
				return
			case val := <-data:
				b.Insert(NewValue(val))
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(b.emitTimer):
				emitted := b.Emit()
				if emitted != nil {
					for _, data := range emitted {
						select {
						case <-done:
							return
						case buffered <- data.GetValue():
						}
					}
				}
			}
		}
	}()

	return buffered
}
