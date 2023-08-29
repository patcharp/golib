package hashing

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// ConvertToByte takes a string and converts it to an array of bytes
func ConvertToByte(payload string) []byte {
	return []byte(payload)
}

// BytesPadding adds padding to the bytes
func BytesPadding(b []byte) []byte {
	const padding = 0
	if len(b) < 4096 {
		for len(b) <= 4096 {
			b = append(b, byte(padding))
		}
	}
	return b
}

// Hash computes a SHA256 hash from a given value
func Hash(value string) string {
	bytes := ConvertToByte(value)
	paddedBytes := BytesPadding(bytes)
	h := sha256.New()
	h.Write(paddedBytes)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func CompareHash(hashed string, value string) bool {
	return hashed == Hash(value)
}

func HashSHA1(value []byte) string {
	h := sha1.New()
	h.Write(value)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func CompareHexSHA1(a, b string) bool {
	aByte, err := hex.DecodeString(a)
	if err != nil {
		return false
	}
	bByte, err := hex.DecodeString(b)
	if err != nil {
		return false
	}
	return CompareSHA1(aByte, bByte)
}

func CompareSHA1(a, b []byte) bool {
	return bytes.Equal(a, b)
}
