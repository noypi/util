package util_test

import (
	"testing"

	"github.com/noypi/util"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestDigiPrint(t *testing.T) {
	assert := assertpkg.New(t)
	s := "abcdefgh"
	const expected = "ab:cd:ef:gh"
	assert.Equal(expected, util.DigiPrint(s))
}
