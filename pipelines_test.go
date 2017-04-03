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
		answer = (Pipelines{pipe1, pipe2}).Run(ch, 1).(int)
		wg.Done()
	}()

	ch <- 3
	close(ch)

	wg.Wait()
	assert.Equal(answer, 8)
}

func TestPipelinesAsyncOut(t *testing.T) {
	assert := assertpkg.New(t)

	in := make(chan interface{}, 1)
	out := make(chan interface{}, 1)

	cleanup := (Pipelines{pipe1, pipe2}).RunAsyncOut(in, out, 1)
	in <- 3
	close(in)

	answer := <-out
	assert.Equal(8, answer.(int))
	close(out)

	cleanup()
}
