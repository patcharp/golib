package attendant

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patcharp/golib/requests"
	"net/http"
	"reflect"
)

// GetAccountByID Get Account by ID API
func (c *Client) GetAccountByID(accountID string) (AccountProfile, error) {
	var account AccountProfile
	data, err := c.get(fmt.Sprintf("/accounts/%s", accountID))
	if err != nil {
		return account, err
	}
	dataByte, err := json.Marshal(data)
	if err := json.Unmarshal(dataByte, &account); err != nil {
		return account, err
	}
	return account, nil
}

func (c *Client) GetAvatarByID(accountID string) (string, []byte, error) {
	var data []byte
	headers := map[string]string{
		echo.HeaderContentType:   "application/json",
		echo.HeaderAuthorization: fmt.Sprintf("%s %s", c.tokenType, c.token),
	}
	uri := fmt.Sprintf("/accounts/%s/avatar", accountID)
	r, err := requests.Get(c.url(uri), headers, bytes.NewBuffer(data), 10)
	if err != nil {
		return "", nil, err
	}
	if r.Code != http.StatusOK {
		return "", nil, errors.New(fmt.Sprintf("server return code %d %s", r.Code, string(r.Body)))
	}
	return r.Header["Content-Type"][0], r.Body, nil
}

// GetTeamMemberByID Get Team Member by ID API
func (c *Client) GetTeamMemberByID(accountID string) ([]CompanyDetail, error) {
	var company []CompanyDetail
	data, err := c.get(fmt.Sprintf("/accounts/%s/teammember", accountID))
	if err != nil {
		return company, err
	}
	dataByte, err := json.Marshal(data)
	if err := json.Unmarshal(dataByte, &company); err != nil {
		return company, err
	}
	return company, nil
}

// GetHeadByID Get Head by ID API
func (c *Client) GetHeadByID(accountID string) ([]CompanyDetail, error) {
	var company []CompanyDetail
	data, err := c.get(fmt.Sprintf("/accounts/%s/head", accountID))
	if err != nil {
		return company, err
	}
	dataByte, err := json.Marshal(data)
	if err := json.Unmarshal(dataByte, &company); err != nil {
		return company, err
	}
	return company, nil
}

// url Set URL Path
func (c *Client) url(path string) string {
	return fmt.Sprintf("%s%s", c.apiEndpoint, path)
}

// get Get Data
func (c *Client) get(uri string) (interface{}, error) {
	var data []byte
	headers := map[string]string{
		echo.HeaderContentType:   "application/json",
		echo.HeaderAuthorization: fmt.Sprintf("%s %s", c.tokenType, c.token),
	}
	r, err := requests.Get(c.url(uri), headers, bytes.NewBuffer(data), 30)
	if err != nil {
		return nil, err
	}
	if r.Code != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("server return code %d %s", r.Code, string(r.Body)))
	}
	var attendantAPIResult APIResult
	if err := json.Unmarshal(r.Body, &attendantAPIResult); err != nil {
		return nil, err
	}
	return reflect.ValueOf(attendantAPIResult.Data).Interface(), nil
}
