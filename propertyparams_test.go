package util_test

import (
	"bytes"
	"testing"

	"github.com/noypi/util"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestParams(t *testing.T) {
	assert := assertpkg.New(t)

	const config = `
somek0 = somev0
# somek1 = comment
  somek1 =  somev1
# somek1 = comment 2
somek2 =   somev 2 "  = some    
somek3 = 35
somek4 = "  got some space "
somek5 = ""  got some space5 ""`

	buf := bytes.NewBufferString(config)
	o, err := util.ReadPropertyParams(buf)
	assert.Nil(err)
	assert.Equal(6, len(o))
	assert.Equal("default value", o.Get("unknown", "default value"))
	assert.Equal(123, o.GetInt("unknown", 123))
	assert.Equal(123.456, o.GetFloat64("unknown", 123.456))
	assert.Equal("somev0", o.Get("somek0", ""))
	assert.Equal("somev1", o.Get("somek1", ""))
	assert.Equal("somev 2 \"  = some", o.Get("somek2", ""))
	assert.Equal("35", o.Get("somek3", ""))
	assert.Equal(35, o.GetInt("somek3", 0))
	assert.Equal("  got some space ", o.Get("somek4", ""))
	assert.Equal("\"  got some space5 \"", o.Get("somek5", ""))

}
