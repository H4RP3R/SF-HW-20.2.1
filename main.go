package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var ErrInvalidInput = errors.New("invalid input: please enter a number")

type Stage func(done chan struct{}, inChan <-chan int) <-chan int

type pipeline struct {
	stages []Stage
}

func (p *pipeline) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

func (p *pipeline) Run(done chan struct{}, dataSource <-chan int) <-chan int {
	c := dataSource
	for _, stage := range p.stages {
		c = stage(done, c)
	}

	return c
}

func NewPipeLine() *pipeline {
	return &pipeline{}
}

func readInput(done chan struct{}) <-chan int {
	outChan := make(chan int)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		defer close(outChan)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Println(err)
			}
			input = strings.TrimSpace(input)
			num, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println(ErrInvalidInput)
			}

			select {
			case outChan <- num:
			case <-done:
				return
			}
		}
	}()

	return outChan
}

func waitForInterrupt(done chan struct{}) {
	var wg sync.WaitGroup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Press Ctrl+C to exit...")
		fmt.Printf("buffer size: %d, delay: %v\n", bufferSize, bufferDelay)
		<-sigChan
		close(done)
		fmt.Println("\nBye!")
	}()

	wg.Wait()
}

func display(done chan struct{}, products <-chan int) {
	go func() {
		for {
			select {
			case num := <-products:
				fmt.Printf("processed: %d\n", num)
			case <-done:
				return
			}
		}
	}()
}

func main() {
	flag.DurationVar(&bufferDelay, "delay", 5*time.Second, "buffer delay")
	flag.IntVar(&bufferSize, "size", 24, "buffer size")
	flag.Parse()

	done := make(chan struct{})
	p := NewPipeLine()
	p.AddStage(filterMultiplesOfThree)
	p.AddStage(filterNegativeNumbers)
	p.AddStage(buffering)

	dataSource := readInput(done)
	products := p.Run(done, dataSource)
	display(done, products)

	waitForInterrupt(done)
}
