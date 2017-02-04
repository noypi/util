package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

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

func (hisPubK *PubKey) DigiPrint() string {
	bb, _ := hisPubK.Marshal()
	return DigiPrint(ToHex(Sha256(bb)))
}

func (hisPubK *PubKey) Verify(msg, signature []byte) error {
	h := sha256.New()
	h.Write(msg)
	return rsa.VerifyPSS((*rsa.PublicKey)(hisPubK), crypto.SHA256, h.Sum(nil), signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
}

type SignedMessage struct {
	Message   []byte
	Signature []byte
}

func (mypriv *PrivKey) SignMessage(bb []byte) (msg *SignedMessage, err error) {
	msg = new(SignedMessage)
	msg.Message = bb
	msg.Signature, err = mypriv.Sign(bb)
	return
}

func (mypriv *PrivKey) SignMessageAndMarshal(bb []byte) ([]byte, error) {
	msg, err := mypriv.SignMessage(bb)
	if nil != err {
		return nil, err
	}
	return SerializeGob(msg)
}

func (hisPubK *PubKey) VerifyMessage(msg *SignedMessage) (err error) {
	return hisPubK.Verify(msg.Message, msg.Signature)
}

func (hisPubK *PubKey) VerifyMessageRaw(bb []byte) (msg []byte, err error) {
	var signedmsg SignedMessage
	if err = DeserializeGob(&signedmsg, bb); nil != err {
		return
	}
	if err = hisPubK.VerifyMessage(&signedmsg); nil != err {
		return
	}
	msg = signedmsg.Message
	return
}
