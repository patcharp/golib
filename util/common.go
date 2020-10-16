package util

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func IsProduction() bool {
	return GetEnv("SERVER_MODE", "dev") == "prod"
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

// Exit program with print elapsed time
func ExitWithCode(startTime time.Time, code int) {
	log.Infoln("Elapsed time", time.Since(startTime).Seconds(), "second(s).")
	os.Exit(code)
}
