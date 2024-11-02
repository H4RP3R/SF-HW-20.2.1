package main

import (
	"fmt"
	"time"

	ringBuf "github.com/H4RP3R/ring_buffer"
)

var (
	bufferDelay time.Duration = 1
	bufferSize  int           = 10
)

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

	ticker := time.NewTicker(bufferDelay * time.Second)
	go func() {
		defer close(outChan)
		for {
			select {
			case <-ticker.C:
				if num, ok := buffer.Pop(); ok {
					outChan <- num
				}
			case <-done:
				return
			}
		}
	}()

	return outChan
}
