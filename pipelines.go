/*
	See also, https://github.com/noypi/fp
*/
package util

type Pipeline func(in, out chan interface{})

type Pipelines []Pipeline

func (pipes Pipelines) Run(in chan interface{}, size int) (out interface{}) {
	if 0 >= size {
		size = 1
	}
	chs := append([]chan interface{}{in}, make([]chan interface{}, len(pipes))...)
	for i := 1; i < len(chs); i++ {
		chs[i] = make(chan interface{}, size)
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

func pipeline_forward(in, out chan interface{}) {
	for v := range in {
		out <- v
	}
}

func (pipes Pipelines) RunAsyncOut(in, out chan interface{}, size int) (cleanup func()) {
	if 0 >= size {
		size = 1
	}

	// create temp channels
	chs := append([]chan interface{}{in}, make([]chan interface{}, len(pipes))...)
	for i := 1; i < len(chs); i++ {
		chs[i] = make(chan interface{}, size)
	}

	for i, j := 0, 1; j < len(chs); i, j = i+1, j+1 {
		go pipes[i](chs[i], chs[j])
	}

	go pipeline_forward(chs[len(chs)-1], out)

	return func() {
		//remove out
		chs = chs[:len(chs)-1]
		//exclude in
		for _, ch := range chs[1:] {
			close(ch)
		}
	}

}
