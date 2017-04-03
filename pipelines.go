/*
	See also, https://github.com/noypi/fp
*/
package util

type Pipeline func(in chan interface{}, out chan interface{})

type Pipelines []Pipeline

func Run(pipes Pipelines, in chan interface{}) (out interface{}) {
	chs := append([]chan interface{}{in}, make([]chan interface{}, len(pipes))...)
	for i := 1; i < len(chs); i++ {
		chs[i] = make(chan interface{}, 1)
	}

	for i, j := 0, 1; j < len(chs); i, j = i+1, j+1 {
		go pipes[i](chs[i], chs[j])
	}

	out = <-chs[len(chs)-1]

	for _, ch := range chs[1:] {
		close(ch)
	}

	return
}
