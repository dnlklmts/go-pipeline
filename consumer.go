package main

import "log"

func consumer(done <-chan bool, src <-chan int) {
	for {
		select {
		case val := <-src:
			log.Printf("Received from pipeline: %d\n", val)
		case <-done:
			return
		}
	}
}
