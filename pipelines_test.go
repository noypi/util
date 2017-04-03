package util

import (
	"sync"
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func pipe1(in chan interface{}, out chan interface{}) {
	for v := range in {
		out <- 2 * v.(int)
	}
}

func pipe2(in chan interface{}, out chan interface{}) {
	for v := range in {
		out <- 2 + v.(int)
	}
}

func TestPipelines(t *testing.T) {
	assert := assertpkg.New(t)

	ch := make(chan interface{}, 1)

	var answer int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		answer = Run(Pipelines{pipe1, pipe2}, ch).(int)
		wg.Done()
	}()

	ch <- 3
	close(ch)

	wg.Wait()
	assert.Equal(answer, 8)
}
