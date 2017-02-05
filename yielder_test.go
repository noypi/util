package util

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func TestYeilderArr(t *testing.T) {
	assert := assertpkg.New(t)

	ns := []interface{}{0, 1, 2, 3, 4, 5}
	fn := YielderArr(ns)
	for i := 0; i < len(ns); i++ {
		b, ok := fn()
		assert.True(ok)
		assert.Equal(ns[i], b)
	}
	b, ok := fn()
	assert.False(ok)
	assert.Nil(b)
}
