package crypto

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

// Random function need to generate random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenSecretString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TODO: How to use
func ValidateString(p string, v []func(rune) bool) bool {
	for _, testRune := range v {
		found := false
		for _, r := range p {
			if testRune(r) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
