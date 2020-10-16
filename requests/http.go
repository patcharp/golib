package requests

import (
	"net/http"
)

const (
	Timeout = 5
)

type Response struct {
	Code   int
	Status string
	Header http.Header
	Body   []byte
}
