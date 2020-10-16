package crypto

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

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

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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
