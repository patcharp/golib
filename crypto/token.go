package crypto

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"strings"
)

func EncodeJWTAccessToken(claims jwt.Claims, signKey *rsa.PrivateKey) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = claims
	return t.SignedString(signKey)
}

func DecodeJWTAccessToken(tokenString string, claims jwt.Claims, signKey *rsa.PrivateKey) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return &signKey.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return token.Claims, nil
	}
	return nil, errors.New("decode token error")
}

func EncryptToken(data []byte, key *rsa.PrivateKey) (string, error) {
	cipherByte, err := RSAEncrypt(key, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherByte), nil
}

func EncryptTokenWithSign(data []byte, key *rsa.PrivateKey) (string, error) {
	cipherByte, err := RSAEncrypt(key, data)
	if err != nil {
		return "", err
	}
	signature, err := SignByRSAKey(key, data)
	token := fmt.Sprintf(
		"%s.%s",
		base64.StdEncoding.EncodeToString(cipherByte),
		base64.StdEncoding.EncodeToString(signature),
	)
	return token, nil
}

func DecryptToken(tokenString string, key *rsa.PrivateKey) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, errors.New("invalid token body")
	}
	plainData, err := RSADecrypt(key, data)
	if err != nil {
		return nil, err
	}
	return plainData, nil
}

func DecryptTokenWithVerifySign(tokenString string, key *rsa.PrivateKey) ([]byte, error) {
	cipherText := strings.Split(tokenString, ".")
	if len(cipherText) != 2 {
		return nil, errors.New("invalid refresh token format")
	}
	data, err := base64.StdEncoding.DecodeString(cipherText[0])
	if err != nil {
		return nil, errors.New("invalid token body")
	}
	plainData, err := RSADecrypt(key, data)
	if err != nil {
		return nil, err
	}
	signature, err := base64.StdEncoding.DecodeString(cipherText[1])
	if err != nil {
		return nil, errors.New("invalid token signature")
	}
	if err := VerifySignedByRSAKey(key.PublicKey, plainData, signature); err != nil {
		return nil, err
	}
	return plainData, nil
}
