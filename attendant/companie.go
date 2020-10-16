package attendant

import (
	"encoding/json"
	"fmt"
	"github.com/patcharp/golib/log"
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
		log.Errorln(pkgName, err, "Json unmarshall team member error")
		return company, nil
	}
	return company, nil
}
