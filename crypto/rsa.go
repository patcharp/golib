package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func InitRSAKey(keySize int) (*rsa.PrivateKey, error) {
	if keySize < 2048 {
		return nil, errors.New("key size is too short")
	}
	return rsa.GenerateKey(rand.Reader, keySize)
}

func ReadRSAPrivateKeyFromPEM(keyBytes []byte, password []byte) (*rsa.PrivateKey, error) {
	var err error
	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		blockBytes, err = x509.DecryptPEMBlock(block, password)
		if err != nil {
			return nil, err
		}
	}
	return x509.ParsePKCS1PrivateKey(blockBytes)
}

func ReadRSAPublicKeyFromPEM(keyBytes []byte, password []byte) (*rsa.PublicKey, error) {
	var err error
	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		blockBytes, err = x509.DecryptPEMBlock(block, password)
		if err != nil {
			return nil, err
		}
	}
	return x509.ParsePKCS1PublicKey(blockBytes)
}

func RSAPrivateKeyToPEM(key *rsa.PrivateKey) string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(key),
			},
		),
	)
}

func RSAPublicKeyToPEM(key *rsa.PublicKey) string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PublicKey(key),
			},
		),
	)
}

func RSAEncrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha512.New(), rand.Reader, &key.PublicKey, data, []byte(""))
}

func RSADecrypt(key *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha512.New(), rand.Reader, key, cipherText, []byte(""))
}

func SignByRSAKey(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	newHash := crypto.SHA512
	pssHash := newHash.New()
	pssHash.Write(data)
	return rsa.SignPSS(
		rand.Reader,
		key,
		crypto.SHA512,
		pssHash.Sum(nil),
		&opts,
	)
}

func VerifySignedByRSAKey(key *rsa.PublicKey, digest []byte, signature []byte) error {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	return rsa.VerifyPSS(
		key,
		crypto.SHA512,
		digest,
		signature,
		&opts,
	)
}
