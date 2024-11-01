package main

import (
	"reflect"
	"testing"
)

func TestFilterNegativeNumbers(t *testing.T) {
	testCases := []struct {
		name        string
		testNumbers []int
		want        []int
	}{
		{
			name:        "No negative numbers",
			testNumbers: []int{18, 2, 0, 35, 4, 100, 0},
			want:        []int{18, 2, 0, 35, 4, 100, 0},
		},
		{
			name:        "Negative numbers",
			testNumbers: []int{-9, 2, -1, 0, 35, -3, 4, 100, 0},
			want:        []int{2, 0, 35, 4, 100, 0},
		},
		{
			name:        "All negative numbers",
			testNumbers: []int{-9, -2, -1, -35, -3, -4, -100},
			want:        []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			done := make(chan struct{})
			inChan := make(chan int)

			go func() {
				defer close(inChan)
				for _, num := range tc.testNumbers {
					inChan <- num
				}
			}()

			outChan := filterNegativeNumbers(done, inChan)
			got := []int{}
			for num := range outChan {
				got = append(got, num)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want: %v, got %v", tc.want, got)
			}
		})
	}
}

func TestFilterMultiplesOfThree(t *testing.T) {
	testCases := []struct {
		name        string
		testNumbers []int
		want        []int
	}{
		{
			name:        "No multiples of 3",
			testNumbers: []int{1, 2, -4, 5, 7, 8, 11, -44},
			want:        []int{},
		},
		{
			name:        "Multiples of 3",
			testNumbers: []int{-1, 0, 3, -12, 55, 15, -6, 31, 99, 77},
			want:        []int{3, -12, 15, -6, 99},
		},
		{
			name:        "All multiples of 3",
			testNumbers: []int{33, -3, 12, 6, -9, -27, 21},
			want:        []int{33, -3, 12, 6, -9, -27, 21},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			done := make(chan struct{})
			inChan := make(chan int)

			go func() {
				defer close(inChan)
				for _, num := range tc.testNumbers {
					inChan <- num
				}
			}()

			outChan := filterMultiplesOfThree(done, inChan)
			got := []int{}
			for num := range outChan {
				got = append(got, num)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want: %v, got %v", tc.want, got)
			}
		})
	}
}

func TestPipeline(t *testing.T) {
	testCases := []struct {
		name        string
		testNumbers []int
		want        []int
		got         []int
	}{
		{
			name:        "Only negative numbers",
			testNumbers: []int{-18, -2, -10, -35, -4, -100, -88},
			want:        nil,
		},
		{
			name:        "No positive multiples of 3",
			testNumbers: []int{-7, -2, 17, 23, 74, 100, -12, 13, 31},
			want:        nil,
		},
		{
			name:        "Mixed",
			testNumbers: []int{5, -2, 33, 1, -15, 23, -8, -21, -4, 3, -12, 0, -6, 9, -10, 17, -18},
			want:        []int{33, 3, 9},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPipeLine()
			p.AddStage(filterMultiplesOfThree)
			p.AddStage(filterNegativeNumbers)
			p.AddStage(buffering)
			dataSource := make(chan int)

			go func() {
				for _, num := range tc.testNumbers {
					dataSource <- num
				}
				close(dataSource)
			}()

			products := p.Run(dataSource)
			for prod := range products {
				tc.got = append(tc.got, prod)
			}

			if !reflect.DeepEqual(tc.got, tc.want) {
				t.Errorf("want: %v, got: %v", tc.want, tc.got)
			}
		})
	}
}
