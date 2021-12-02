package requests

import (
	"net/http"
)

const (
	Timeout = 10
)

type Response struct {
	Code   int
	Status string
	Header http.Header
	Body   []byte
}
