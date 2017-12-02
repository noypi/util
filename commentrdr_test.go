package util_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/noypi/util"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestCommentRdr(t *testing.T) {
	assert := assertpkg.New(t)

	const expected = ` some file
 ok # have read
davidking mamay
1
2`

	const config = `some
# some file
 start config    # ok # have read
#davidking mamay
got
#1
so
3 #2`

	const config2 = `some
# some file
 start config    # ok # have read


#davidking mamay
got

	
#1
so

3 #2`

	const config3 = `some
# some file
 start config    # ok # have read


#davidking mamay
got

	
#1
so

3 #2
`
	const config4 = `some
# some file
 start config    # ok # have read


#davidking mamay
got

	
#1
so

3 #2


`

	for i, s := range []string{config, config2, config3, config4} {
		buf := bytes.NewBufferString(s)
		r := util.NewCommentReader(buf, '#')
		bb, err := ioutil.ReadAll(r)
		assert.Nil(err)
		if 1 < i {
			assert.Equal(expected+"\n", string(bb))
		} else {
			assert.Equal(expected, string(bb))
		}

	}
}

func TestCommentRdrFileOpts(t *testing.T) {
	assert := assertpkg.New(t)

	content := `#
# some opt
#    my opt = value of opt 1  
# some word
# my Opt 2 =     some opt 2

# my opt 3 = "  is not included "
	
`
	buf := bytes.NewBufferString(content)
	opts, err := util.ReadPropertyParams(util.NewCommentReader(buf, '#'))
	assert.Nil(err)

	assert.Equal("value of opt 1", opts.Get("my opt", ""))
	assert.Equal("some opt 2", opts.Get("my opt 2", ""))
	assert.Equal("  is not included ", opts.Get("my opt 3", ""))

}

func TestCommentRdrParseScript(t *testing.T) {
	assert := assertpkg.New(t)

	var script = `# chromosomes = 4
# genes = 3

symbols %s
alias namedind %s min=%.1f max=%.1f

beginscan
	namedind
	linreg.slope namedind min=%.1f max=%.1f
	winif gainafter days=3
endscan
`

	buf := bytes.NewBufferString(script)
	opts, err := util.ReadPropertyParams(util.NewCommentReader(buf, '#'))
	assert.Nil(err)

	assert.Equal("3", opts.Get("genes", ""))
	assert.Equal("4", opts.Get("chromosomes", ""))
	assert.Equal(3, opts.GetInt("genes", 0))
	assert.Equal(4, opts.GetInt("chromosomes", 0))

}
