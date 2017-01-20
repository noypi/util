package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// implements default encrypt/decrypt using rsa,
// will not sign

type PrivKey rsa.PrivateKey
type PubKey rsa.PublicKey

func GenPrivKey(bits int) (*PrivKey, error) {
	privk, err := rsa.GenerateKey(rand.Reader, bits)
	return (*PrivKey)(privk), err
}

func (this *PrivKey) Marshal() []byte {
	return x509.MarshalPKCS1PrivateKey((*rsa.PrivateKey)(this))
}

func (this *PrivKey) PubKey() *PubKey {
	return (*PubKey)(&((*rsa.PrivateKey)(this).PublicKey))
}

func (this *PrivKey) MarshalPem() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: this.Marshal(),
	})
}

func (this *PubKey) Marshal() (bbPub []byte, err error) {
	bbPub, err = x509.MarshalPKIXPublicKey((*rsa.PublicKey)(this))
	return
}

func (this *PubKey) MarshalPem() ([]byte, error) {
	bb, err := this.Marshal()
	if nil != err {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: bb,
	}), nil
}

func ParsePublickey(bb []byte) (*PubKey, error) {
	pubKif, err := x509.ParsePKIXPublicKey(bb)
	if nil != err {
		return nil, err
	}
	return (*PubKey)(pubKif.(*rsa.PublicKey)), nil
}

func ParsePublickeyPem(bbPem []byte) (*PubKey, error) {
	block, _ := pem.Decode(bbPem)
	return ParsePublickey(block.Bytes)
}

func (this *PrivKey) Decrypt(bb []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, (*rsa.PrivateKey)(this), bb, []byte(""))
}

func (this *PubKey) Encrypt(bb []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, (*rsa.PublicKey)(this), bb, []byte(""))
}

func (mypriv *PrivKey) Sign(bb []byte) (signature []byte, err error) {
	h := sha256.New()
	h.Write(bb)
	signature, err = rsa.SignPSS(rand.Reader, (*rsa.PrivateKey)(mypriv), crypto.SHA256, h.Sum(nil), &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
	return
}

func (hisPubK *PubKey) Verify(msg, signature []byte) error {
	h := sha256.New()
	h.Write(msg)
	return rsa.VerifyPSS((*rsa.PublicKey)(hisPubK), crypto.SHA256, h.Sum(nil), signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
}
