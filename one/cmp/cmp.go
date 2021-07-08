package cmp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/patcharp/golib/v2/requests"
	"github.com/patcharp/golib/v2/util/httputil"
	"net/http"
)

func New() Client {
	return Client{ApiEndpoint: ProductionEndpoint}
}

func (c *Client) GetConsentStatus(form ConsentRequestForm) (StatusResult, error) {
	var result StatusResult
	body, _ := json.Marshal(&form)
	r, err := c.send(http.MethodPost, c.url("/api/chat/service/check-consent"), body, nil)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (c *Client) send(method string, url string, body []byte, customHeader map[string]string) (requests.Response, error) {
	headers := map[string]string{
		httputil.HeaderContentType: "application/json",
	}
	if customHeader != nil {
		for k, v := range customHeader {
			headers[k] = v
		}
	}
	r, err := requests.Request(method, url, headers, bytes.NewBuffer(body), 0)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (c *Client) url(path string) string {
	return fmt.Sprintf("%s%s", c.ApiEndpoint, path)
}
