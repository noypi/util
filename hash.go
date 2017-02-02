package util

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

func SplitSep2(s, sep string) string {
	buf := bytes.NewBufferString("")
	buf.Grow(len(s) + (len(s) / 2))
	for i, c := range s {
		if 0 != i && 0 == (i&1) {
			buf.WriteRune(':')
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

func ToHex(bb []byte) string {
	return fmt.Sprintf("%x", bb)
}

func Sha256(bb []byte) []byte {
	return Hash(bb, sha256.New())
}

func Sha512(bb []byte, h hash.Hash) []byte {
	return Hash(bb, sha512.New())
}

func Hash(bb []byte, h hash.Hash) []byte {
	h.Write(bb)
	return h.Sum(nil)
}
