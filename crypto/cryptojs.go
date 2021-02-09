package crypto

import (
	"github.com/Luzifer/go-openssl/v4"
)

func CryptoJSEncrypt(data []byte, key string) ([]byte, error) {
	o := openssl.New()
	return o.EncryptBytes(key, data, openssl.BytesToKeyMD5)
}

func CryptoJSDecrypt(cipher []byte, key string) ([]byte, error) {
	o := openssl.New()
	return o.DecryptBytes(key, cipher, openssl.BytesToKeyMD5)
}
