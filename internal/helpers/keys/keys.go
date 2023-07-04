package keys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/amleonc/tabula/config"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// ------------ Variables ------------ //

var (
	privateKey = readPrivateKey()
	publicKey  = readPublicKey()

	pvk = readJWKPrivateKey()
	pbk = readJWKPublicKey()
)

// ------------ Functions ------------ //

func PrivateRSAKey() *rsa.PrivateKey {
	return privateKey
}

func PublicRSAKey() *rsa.PublicKey {
	return publicKey
}

func PrivateJWKKey() jwk.Key {
	return pvk
}

func PublicJWKKey() jwk.Key {
	return pbk
}

func readPrivateKey() *rsa.PrivateKey {
	keyBytes, err := os.ReadFile(config.PrivateKeyPath())
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatal("cannot decode key file contents")
	}
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	return rsaPrivateKey
}

func readPublicKey() *rsa.PublicKey {
	return &privateKey.PublicKey
}

func readJWKPrivateKey() jwk.Key {
	k, err := jwk.FromRaw(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	return k
}

func readJWKPublicKey() jwk.Key {
	k, err := jwk.FromRaw(publicKey)
	if err != nil {
		log.Fatal(err)
	}
	return k
}
