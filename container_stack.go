package util

import (
	"container/list"
)

type Stack list.List

func NewStack() *Stack {
	return (*Stack)(list.New())
}

func (this *Stack) Push(e interface{}) {
	(*list.List)(this).PushBack(e)
}

func (this *Stack) Pop() interface{} {
	e := (*list.List)(this).Back()
	(*list.List)(this).Remove(e)
	return e.Value
}

func (this *Stack) Len() int {
	return (*list.List)(this).Len()
}
