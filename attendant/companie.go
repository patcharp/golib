package attendant

import (
	"encoding/json"
	"fmt"
)

// GetCompanies Get All Companies API
func (c *Client) GetCompanies() ([]Companies, error) {
	var company []Companies
	data, err := c.get(fmt.Sprintf("/companies"))
	if err != nil {
		return company, err
	}
	dataByte, err := json.Marshal(data)
	if err := json.Unmarshal(dataByte, &company); err != nil {
		return company, err
	}
	return company, nil
}

// GetCompanyEmployee
func (c *Client) GetCompanyEmployee(taxNo string) ([]StaffInfo, error) {
	var employees []StaffInfo
	data, err := c.get(fmt.Sprintf("/companies/%s/employees", taxNo))
	if err != nil {
		return nil, err
	}
	dataByte, err := json.Marshal(data)
	if err := json.Unmarshal(dataByte, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}
