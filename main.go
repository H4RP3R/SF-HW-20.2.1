package main

import "fmt"

type Stage func(done chan struct{}, inChan <-chan int) <-chan int

type pipeline struct {
	stages []Stage
	done   chan struct{}
}

func (p *pipeline) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

func (p *pipeline) Run(dataSource <-chan int) <-chan int {
	c := dataSource

	for _, stage := range p.stages {
		c = stage(p.done, c)
	}

	return c
}

func NewPipeLine() *pipeline {
	return &pipeline{
		done: make(chan struct{}),
	}
}

func main() {
	p := NewPipeLine()
	p.AddStage(filterMultiplesOfThree)
	p.AddStage(filterNegativeNumbers)
	p.AddStage(buffering)
	dataSource := make(chan int)

	numbers := []int{12, -12, 21, 2, 3, 33, 4, -5, -7, -99, 13, 15, 1}

	go func() {
		for _, n := range numbers {
			dataSource <- n
		}
		close(dataSource)
	}()

	products := p.Run(dataSource)

	for p := range products {
		fmt.Print(p, " ")
	}
}
