package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

type PrivKey rsa.PrivateKey

func GenPrivKey(bits int) (*PrivKey, error) {
	privk, err := rsa.GenerateKey(rand.Reader, bits)
	return (*PrivKey)(privk), err
}

func (this *PrivKey) Marshal() []byte {
	return x509.MarshalPKCS1PrivateKey((*rsa.PrivateKey)(this))
}

func (this *PrivKey) MarshalPem() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: this.Marshal(),
	})
}

func (this *PrivKey) MarshalPublicKey() (bbPub []byte, err error) {
	pubk := &((*rsa.PrivateKey)(this).PublicKey)
	bbPub, err = x509.MarshalPKIXPublicKey(pubk)
	return
}

func (this PrivKey) MarshalPublicKeyPem() ([]byte, error) {
	bb, err := this.MarshalPublicKey()
	if nil != err {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: bb,
	}), nil
}

func ParsePublickeyPem(bbPem []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(bbPem)
	pubKif, err := x509.ParsePKIXPublicKey(block.Bytes)
	if nil != err {
		return nil, err
	}
	return pubKif.(*rsa.PublicKey), nil
}

func EncryptUsingPubKey(bb []byte, pubk *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pubk, bb, []byte(""))
}
