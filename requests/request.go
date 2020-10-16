package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

func Request(method string, url string, headers map[string]string, body io.Reader, timeout int) (Response, error) {
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

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
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
