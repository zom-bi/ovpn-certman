package handlers

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"git.klink.asia/paul/certman/views"

	"github.com/jinzhu/gorm"
)

func ListCertHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := views.New(req)
		v.Render(w, "cert_list")
	}
}

func GenCertHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := views.New(req)

		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("Could not generate keypair: %s", err)
		}

		caCert, caKey, err := loadX509KeyPair("ca.crt", "ca.key")
		if err != nil {
			v.Render(w, "500")
			log.Fatalf("error loading ca keyfiles: %s", err)
		}

		derBytes, err := CreateCertificate(key, caCert, caKey)
		pem.Encode(w, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

		pkBytes := x509.MarshalPKCS1PrivateKey(key)
		pem.Encode(w, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: pkBytes})
		return
	}
}

func loadX509KeyPair(certFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	cf, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}

	kf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	cpb, cr := pem.Decode(cf)
	fmt.Println(string(cr))
	kpb, kr := pem.Decode(kf)
	fmt.Println(string(kr))
	crt, err := x509.ParseCertificate(cpb.Bytes)

	if err != nil {
		return nil, nil, err
	}
	key, err := x509.ParsePKCS1PrivateKey(kpb.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return crt, key, nil
}

// CreateCertificate creates a CA-signed certificate
func CreateCertificate(key interface{}, caCert *x509.Certificate, caKey interface{}) ([]byte, error) {
	subj := caCert.Subject
	// .. except for the common name
	subj.CommonName = "clientName"

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Obscure error in cert serial number generation: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subj,

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour * 356 * 5),

		SignatureAlgorithm: x509.SHA256WithRSA,
		//KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	return x509.CreateCertificate(rand.Reader, &template, caCert, publicKey(key), caKey)
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}
