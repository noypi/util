package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
	"time"
)

// comma separted hosts / IP
func GenerateCert(org, hosts string, bits int) (key, cert []byte, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)

	tValidFrom := time.Now().AddDate(-1, 0, 0)
	tValidTo := tValidFrom.AddDate(2, 0, 0)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if nil != err {
		return
	}
	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{org},
		},
		NotBefore:             tValidFrom,
		NotAfter:              tValidTo,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}
	for _, h := range strings.Split(hosts, ",") {
		if ip := net.ParseIP(h); nil != ip {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, h)
		}
	}

	bb, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if nil != err {
		return
	}

	certbuf := bytes.NewBufferString("")
	pem.Encode(certbuf, &pem.Block{Type: "CERTIFICATE", Bytes: bb})

	keybuf := bytes.NewBufferString("")
	pem.Encode(keybuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return keybuf.Bytes(), certbuf.Bytes(), nil
}
