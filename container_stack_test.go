package util

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	assert := assertpkg.New(t)
	stack := NewStack()

	assert.Equal(0, stack.Len())
	stack.Push(12)
	stack.Push(24)
	assert.Equal(2, stack.Len())
	stack.Push("25")
	assert.Equal(3, stack.Len())

	assert.Equal("25", stack.Pop())
	assert.Equal(2, stack.Len())
	assert.Equal(24, stack.Pop())
	assert.Equal(1, stack.Len())
	assert.Equal(12, stack.Pop())
	assert.Equal(0, stack.Len())
}
