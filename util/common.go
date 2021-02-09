package util

import (
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func AtoI(s string, v int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return v
	}
	return i
}

func AtoF(s string, v float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return v
	}
	return f
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

// Real display width string - remove floating char
func StringDisplayWidth(str string) int {
	floatingChar := map[string][]string{
		"thai": []string{"ิ", "ี", "ึ", "ื", "ุ", "ู", "่", "้", "๊", "๋", "์"},
	}
	count := 0
	for _, l := range floatingChar {
		for _, v := range l {
			count += strings.Count(str, v)
		}
	}
	return utf8.RuneCountInString(str) - count
}
