package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

var DefaultTLSConfig = tls.Config{InsecureSkipVerify: true}

func Request(method string, url string, headers map[string]string, body io.Reader, timeout int) (Response, error) {
	return RequestWithTLSConfig(method, url, headers, body, timeout, nil)
}

func RequestWithTLSConfig(method string, url string, headers map[string]string, body io.Reader, timeout int, tlsCfg *tls.Config) (Response, error) {
	if timeout == 0 {
		timeout = Timeout
	}
	r := Response{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return r, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if tlsCfg == nil {
		tlsCfg = &DefaultTLSConfig
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsCfg
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, err
	}

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)

	r.Code = resp.StatusCode
	r.Status = resp.Status
	r.Body = buf.Bytes()
	r.Header = resp.Header
	_ = resp.Body.Close()

	return r, nil
}

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
