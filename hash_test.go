package util_test

import (
	"testing"

	"github.com/noypi/util"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestSplitSep(t *testing.T) {
	assert := assertpkg.New(t)
	s := "abcdefgh"
	const expected = "ab:cd:ef:gh"
	assert.Equal(expected, util.SplitSep2(s, ":"))
}
