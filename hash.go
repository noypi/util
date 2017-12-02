package util

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

func DigiPrint(s string) string {
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

func Sha256(bb []byte) []byte {
	return Hash(bb, sha256.New())
}

func Sha512(bb []byte) []byte {
	return Hash(bb, sha512.New())
}

func Hash(bb []byte, h hash.Hash) []byte {
	h.Write(bb)
	return h.Sum(nil)
}

func SaltedHash256(data, salt []byte) (hash, outSalt []byte) {
	return Sha256(append(data, salt...)), salt
}

func SaltedHash512(data, salt []byte) (hash, outSalt []byte) {
	return Sha512(append(data, salt...)), salt
}

func GenSalt(size int) []byte {
	bbSalt := make([]byte, size)
	rand.Read(bbSalt)
	return bbSalt
}
