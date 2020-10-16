package imagik

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"time"
)

const MaxTimeout = 10

func UrlGrabber(url string, headers map[string]string, b *[]byte, mimeType *string, timeout int) error {
	if timeout <= 0 || timeout >= MaxTimeout {
		timeout = MaxTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	var data io.Reader
	req, err := http.NewRequest(http.MethodGet, url, data)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	*b = buf.Bytes()
	*mimeType = resp.Header.Get(echo.HeaderContentType)
	return resp.Body.Close()
}
