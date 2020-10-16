package requests

import (
	"io"
	"net/http"
)

func Get(url string, headers map[string]string, body io.Reader, timeout int) (Response, error) {
	return Request(http.MethodGet, url, headers, body, timeout)
}

func Post(url string, headers map[string]string, body io.Reader, timeout int) (Response, error) {
	return Request(http.MethodPost, url, headers, body, timeout)
}

func Put(url string, headers map[string]string, body io.Reader, timeout int) (Response, error) {
	return Request(http.MethodPut, url, headers, body, timeout)
}

func Delete(url string, headers map[string]string, body io.Reader, timeout int) (Response, error) {
	return Request(http.MethodDelete, url, headers, body, timeout)
}
