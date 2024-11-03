package main

import (
	"fmt"
	"time"

	ringBuf "github.com/H4RP3R/ring_buffer"
)

var (
	bufferDelay time.Duration
	bufferSize  int
)

// filterNegativeNumbers returns a channel that emits only non-negative
// numbers from the input channel.
func filterNegativeNumbers(done chan struct{}, inChan <-chan int) <-chan int {
	outChan := make(chan int)

	go func() {
		defer close(outChan)
		for num := range inChan {
			if num >= 0 {
				select {
				case outChan <- num:
				case <-done:
					return
				}
			}
		}
	}()

	return outChan
}

// filterMultiplesOfThree returns a channel that emits only numbers that
// are multiples of 3.
func filterMultiplesOfThree(done chan struct{}, inChan <-chan int) <-chan int {
	outChan := make(chan int)

	go func() {
		defer close(outChan)
		for num := range inChan {
			if num != 0 && num%3 == 0 {
				select {
				case outChan <- num:
				case <-done:
					return
				}
			}
		}
	}()

	return outChan
}

// buffering collects numbers from the input channel into a buffer and emits
// them to the output channel at intervals specified by bufferDelay.
func buffering(done chan struct{}, inChan <-chan int) <-chan int {
	outChan := make(chan int)
	buffer, err := ringBuf.New[int](bufferSize)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for num := range inChan {
			buffer.Push(num)
		}
	}()

	ticker := time.NewTicker(bufferDelay)
	go func() {
		defer func() {
			ticker.Stop()
			close(outChan)
		}()
		for {
			select {
			case <-ticker.C:
				for !buffer.IsEmpty() {
					if num, ok := buffer.Pop(); ok {
						outChan <- num
					}
				}
			case <-done:
				return
			}
		}
	}()

	return outChan
}
