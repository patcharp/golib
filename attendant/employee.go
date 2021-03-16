package attendant

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patcharp/golib/requests"
	"net/http"
)

// GetEmployeeByID Get employee by ID API
func (c *Client) GetEmployeeByID(employeeId string) (StaffInfo, error) {
	var employee StaffInfo
	data, err := c.get(fmt.Sprintf("/employees/%s", employeeId))
	if err != nil {
		return employee, err
	}
	dataByte, err := json.Marshal(data)
	if err := json.Unmarshal(dataByte, &employee); err != nil {
		return employee, err
	}
	return employee, nil
}

// GetEmployeeAvatar avatar
func (c *Client) GetEmployeeAvatar(employeeId string) (string, []byte, error) {
	var data []byte
	headers := map[string]string{
		echo.HeaderContentType:   "application/json",
		echo.HeaderAuthorization: fmt.Sprintf("%s %s", c.tokenType, c.token),
	}
	uri := fmt.Sprintf("/employees/%s/avatar", employeeId)
	r, err := requests.Get(c.url(uri), headers, bytes.NewBuffer(data), 10)
	if err != nil {
		return "", nil, err
	}
	if r.Code != http.StatusOK {
		return "", nil, errors.New(fmt.Sprintf("server return code %d %s", r.Code, string(r.Body)))
	}
	return r.Header["Content-Type"][0], r.Body, nil
}
