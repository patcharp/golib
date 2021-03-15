package attendant

import (
	"encoding/json"
	"fmt"
)

// GetCompanies Get All Companies API
func (client *Client) GetCompanies() ([]Companies, error) {
	var company []Companies
	data, err := client.get(fmt.Sprintf("/companies"))
	if err != nil {
		return company, err
	}
	dataByte, err := json.Marshal(data)
	if err := json.Unmarshal(dataByte, &company); err != nil {
		return company, err
	}
	return company, nil
}
