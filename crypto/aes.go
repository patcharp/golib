package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func AESEncrypt(data []byte, key string) (string, error) {
	gcmCipher, err := genGcmCipher(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcmCipher.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(gcmCipher.Seal(nonce, nonce, data, nil)), nil
}

func AESDecrypt(cipherText string, key string) ([]byte, error) {
	cipherByte, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	var data []byte
	gcmCipher, err := genGcmCipher(key)
	if err != nil {
		return data, err
	}
	nonceSize := gcmCipher.NonceSize()
	if len(cipherByte) < nonceSize {
		return nil, errors.New("invalid string size")
	}
	nonce, cipherByte := cipherByte[:nonceSize], cipherByte[nonceSize:]
	return gcmCipher.Open(nil, nonce, cipherByte, nil)
}

func genGcmCipher(key string) (cipher.AEAD, error) {
	cipherBlock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcmCipher, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}
	return gcmCipher, nil
}
