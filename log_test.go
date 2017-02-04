package util_test

import (
	"testing"

	"github.com/noypi/util"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	assert := assertpkg.New(t)
	s := util.StackTrace(1)
	assert.Contains(s, "TestLog", "s=%s", s)
	assert.Contains(s, "log_test")
}
