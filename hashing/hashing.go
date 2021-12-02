package hashing

import (
	"crypto/sha256"
	"encoding/hex"
)

type bytes []byte

// ConvertToByte takes a string and converts it to an array of bytes
func ConvertToByte(payload string) []byte {
	return []byte(payload)
}

// BytesPadding adds padding to the bytes
func BytesPadding(b bytes) []byte {
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
