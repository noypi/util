package util

import (
	"bytes"
	"io"
)

type CommentReader struct {
	CommentChar rune
	r           io.Reader
}

func NewCommentReader(r io.Reader, char rune) io.Reader {
	o := new(CommentReader)
	o.CommentChar = char
	o.r = r
	return o
}

func (this *CommentReader) Read(p []byte) (n int, err error) {
	var p2 = make([]byte, len(p))
	var n2 int
	var bInComment bool
	for 0 < len(p2) && n < len(p) && io.EOF != err {
		n2, err = this.r.Read(p2)
		if 0 == n2 {
			break
		}
		var p3 = p2
		for 0 < len(p3) && 0 < n2 {
			// there might be multiple lines here
			var i int
			if !bInComment {
				i = bytes.IndexRune(p3, this.CommentChar)
			}
			if 0 <= i {
				bInComment = true
				i++
				j := i
				n2 -= i
				for ; j < len(p3) && 0 < n2; j++ {
					if '\n' == p3[j] {
						bInComment = false
						p[n] = '\n'
						n++
						break
					}
					p[n] = p3[j]
					n++
					n2--
				}

				if n < len(p3) {
					p3 = p3[j:]
				} else {
					p3 = nil
				}
			} else {
				break
			}

		}
		p2 = p2[0:n]
	}

	return
}
