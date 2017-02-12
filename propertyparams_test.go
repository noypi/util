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

func TestParamsEnvVars(t *testing.T) {
	assert := assertpkg.New(t)

	const config = `
$Var0 = var value
somek0 = $Var0
somek1 = $Var1
$Var1 = var value2
somek2 = $Var1
somek3 =  some $Var0
$Var3 = $Var0 concat
somek4 =  $Var3
somek5 =  ${Var3}concat`

	buf := bytes.NewBufferString(config)
	o, err := util.ReadPropertyParams(buf)
	assert.Nil(err)
	assert.Equal("var value", o.Get("somek0", ""))
	assert.Equal("", o.Get("somek1", "aaaa"))
	assert.Equal("var value2", o.Get("somek2", ""))
	assert.Equal("some var value", o.Get("somek3", ""))
	assert.Equal("var value concat", o.Get("somek4", ""))
	assert.Equal("var value concatconcat", o.Get("somek5", ""))
}

func TestParamsNamespace(t *testing.T) {
	assert := assertpkg.New(t)

	const config = `
$Var0 = var value
somek0 = $Var0
somek1 = $Var1
$Var1 = var value2
$Var3 = $Var0 concat

[ns1]
somek2 = $Var1
somek3 =  some $Var0

[[ns2]]
somek4 =  $Var3
somek5 =  ${Var3}concat

[ns1b]
somek2 = $Var1 b
somek3 =  some $Var0 b

[[ns2b]]
somek4 =  $Var3 b
somek5 =  ${Var3}concat b

[[ns2c]]
somek5 =  ${Var3}concat c

[[[ns3c]]]
somek5 =  ${Var3}concat c

[[ns2d]]
somek5 =  ${Var3}concat d

[ns1e]
somek5 =  ${Var3}concat e`

	buf := bytes.NewBufferString(config)
	o, err := util.ReadPropertyParams(buf)
	assert.Nil(err)
	assert.Equal("var value", o.Get("somek0", ""))
	assert.Equal("", o.Get("somek1", "aaaa"))
	assert.Equal("<notfound>", o.Get("somek2", "<notfound>"))

	assert.Equal("var value2", o.Get("ns1.somek2", ""))
	assert.Equal("some var value", o.Get("ns1.somek3", ""))
	assert.Equal("var value concat", o.Get("ns1.ns2.somek4", ""), "k=%s: m=%v", "ns1.ns2.somek4", o)
	assert.Equal("var value concatconcat", o.Get("ns1.ns2.somek5", ""))

	assert.Equal("var value2 b", o.Get("ns1b.somek2", ""), "k=%s m=%v", "ns1b.somek2", o)
	assert.Equal("some var value b", o.Get("ns1b.somek3", ""))
	assert.Equal("var value concat b", o.Get("ns1b.ns2b.somek4", ""))
	assert.Equal("var value concatconcat b", o.Get("ns1b.ns2b.somek5", ""))
	assert.Equal("var value concatconcat c", o.Get("ns1b.ns2c.somek5", ""))
	assert.Equal("var value concatconcat c", o.Get("ns1b.ns2c.ns3c.somek5", ""))
	assert.Equal("var value concatconcat d", o.Get("ns1b.ns2d.somek5", ""))
	assert.Equal("var value concatconcat e", o.Get("ns1e.somek5", ""))
}

func TestParamsNamespaceCopy(t *testing.T) {
	assert := assertpkg.New(t)

	const config = `
$Var0 = var value
somek0 = $Var0
somek1 = $Var1
$Var1 = var value2
$Var3 = $Var0 concat

[ns1]
somek2 = $Var1
somek3 =  some $Var0

[[ns2]]
somek4 =  $Var3
somek5 =  ${Var3}concat

[ns1b]
somek2 = $Var1 b
somek3 =  some $Var0 b

[[ns2b]]
somek4 =  $Var3 b
somek5 =  ${Var3}concat b

[[ns2c]]
somek5 =  ${Var3}concat c

[[[ns3c]]]
:copy ns1.ns2

[[[ns3d]]]
:copy ns1.ns2
somek5 = overwritten

`

	buf := bytes.NewBufferString(config)
	o, err := util.ReadPropertyParams(buf)
	assert.Nil(err)
	assert.Equal("var value", o.Get("somek0", ""))
	assert.Equal("", o.Get("somek1", "aaaa"))
	assert.Equal("<notfound>", o.Get("somek2", "<notfound>"))

	assert.Equal("var value2", o.Get("ns1.somek2", ""))
	assert.Equal("some var value", o.Get("ns1.somek3", ""))
	assert.Equal("var value concat", o.Get("ns1.ns2.somek4", ""), "k=%s: m=%v", "ns1.ns2.somek4", o)
	assert.Equal("var value concatconcat", o.Get("ns1.ns2.somek5", ""))

	assert.Equal("var value2 b", o.Get("ns1b.somek2", ""), "k=%s m=%v", "ns1b.somek2", o)
	assert.Equal("some var value b", o.Get("ns1b.somek3", ""))
	assert.Equal("var value concat b", o.Get("ns1b.ns2b.somek4", ""))
	assert.Equal("var value concatconcat b", o.Get("ns1b.ns2b.somek5", ""))
	assert.Equal("var value concatconcat c", o.Get("ns1b.ns2c.somek5", ""))

	// copied
	assert.Equal("var value concat", o.Get("ns1b.ns2c.ns3c.somek4", ""), "o=%v", o)
	assert.Equal("var value concatconcat", o.Get("ns1b.ns2c.ns3c.somek5", ""))

	// copied
	assert.Equal("var value concat", o.Get("ns1b.ns2c.ns3d.somek4", ""), "o=%v", o)
	assert.Equal("overwritten", o.Get("ns1b.ns2c.ns3d.somek5", ""))

}
