package identity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patcharp/golib/log"
	"github.com/patcharp/golib/requests"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"reflect"
)

const (
	pkgName               = "IDENTITY"
	OrgProductionEndpoint = "https://one.th/api/v2/service/business"
)

func NewOrgClient(username string, password string, clientId string, clientSecret string, endpoint string) (OrgClient, error) {
	const orgTokenCacheKey = "org_token_key"
	var org OrgClient
	var r AuthenticationResult
	if endpoint == "" {
		endpoint = OrgProductionEndpoint
	}
	org = OrgClient{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		TokenType:    r.TokenType,
		ApiEndpoint:  endpoint,
	}
	return org, nil
}

func (org *OrgClient) GetDepartments(taxNo string) ([]Department, error) {
	var dept []Department
	data, err := org.get("/department", taxNo)
	if err != nil {
		return dept, err
	}
	for _, v := range data.([]interface{}) {
		vMap := v.(map[string]interface{})
		var id uuid.UUID
		parentDeptId := uuid.Nil
		if id = uuid.FromStringOrNil(vMap["id"].(string)); id != uuid.Nil {
			if vMap["parent_dept_id"] != nil {
				parentDeptId = uuid.FromStringOrNil(vMap["parent_dept_id"].(string))
			}
			d := Department{
				Id:           id,
				Name:         vMap["dept_name"].(string),
				ParentDeptId: &parentDeptId,
			}
			dept = append(dept, d)
		} else {
			log.Errorln(pkgName, errors.New(fmt.Sprint("Invalid department", v)), "Get department error")
		}
	}
	return dept, nil
}

func (org *OrgClient) GetAccounts(taxNo string) ([]AccountProfile, error) {
	var accounts []AccountProfile
	data, err := org.get("/account", taxNo)
	if err != nil {
		return accounts, err
	}
	for _, v := range data.([]interface{}) {
		jsonMap, _ := json.Marshal(v.(map[string]interface{}))
		var a AccountProfile
		if err := json.Unmarshal(jsonMap, &a); err == nil {
			accounts = append(accounts, a)
		} else {
			log.Errorln(pkgName, err, "Unmarshal account")
		}
	}
	return accounts, nil
}

func (org *OrgClient) GetDepartmentAccounts(taxNo string, departmentUid uuid.UUID) ([]Employee, error) {
	var employee []Employee
	data, err := org.get(fmt.Sprintf("/department/%s", departmentUid), taxNo)
	if err != nil {
		return employee, err
	}
	// Get position
	positionCode := map[uuid.UUID]string{}
	for _, v := range data.(map[string]interface{})["has_role"].([]interface{}) {
		roleId := uuid.FromStringOrNil(v.(map[string]interface{})["role_id"].(string))
		positionCode[roleId] = v.(map[string]interface{})["role"].(map[string]interface{})["role_name"].(string)
	}
	//
	for _, v := range data.(map[string]interface{})["has_account"].([]interface{}) {
		jsonMap, _ := json.Marshal(v)
		var e Employee
		if err := json.Unmarshal(jsonMap, &e); err == nil {
			e.Position = positionCode[e.PositionId]
			employee = append(employee, e)
		} else {
			log.Errorln(pkgName, err, "Unmarshal employee")
		}
	}
	return employee, nil
}

func (org *OrgClient) SetEndpoint(ep string) {
	org.ApiEndpoint = ep
}

func (org *OrgClient) get(uri string, taxNo string) (interface{}, error) {
	data, _ := json.Marshal(&struct {
		TaxNo string `json:"tax_id"`
	}{
		TaxNo: taxNo,
	})
	headers := map[string]string{
		echo.HeaderContentType:   "application/json",
		echo.HeaderAuthorization: fmt.Sprintf("%s %s", org.TokenType, org.AccessToken),
	}
	r, err := requests.Get(org.url(uri), headers, bytes.NewBuffer(data), 30)
	if err != nil {
		return nil, err
	}
	if r.Code != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("server return code %d %s", r.Code, string(r.Body)))
	}
	var orgApiResult OrgApiResult
	if err := json.Unmarshal(r.Body, &orgApiResult); err != nil {
		return nil, err
	}
	return reflect.ValueOf(orgApiResult.Data).Interface(), nil
}

func (org *OrgClient) url(path string) string {
	return fmt.Sprintf("%s%s", org.ApiEndpoint, path)
}

func (org *OrgClient) GetSubordinateDepartmentAccounts(accountId string, taxNo string) ([]TeamMember, error) {
	var teamMembers []TeamMember
	rawData, err := org.get(fmt.Sprintf("/account/%s/subordinate-department", accountId), taxNo)
	if err != nil {
		return teamMembers, err
	}
	dataByte, err := json.Marshal(rawData)
	if err := json.Unmarshal(dataByte, &teamMembers); err != nil {
		log.Errorln(pkgName, err, "Json unmarshall team member error")
		return teamMembers, nil
	}
	return teamMembers, nil
}

func (org *OrgClient) GetHeadDepartmentAccounts(accountId string, taxNo string) ([]HeadDepartment, error) {
	var headDepart []HeadDepartment
	rawData, err := org.get(fmt.Sprintf("/account/%s/head-department", accountId), taxNo)
	if err != nil {
		return headDepart, err
	}
	dataByte, err := json.Marshal(rawData)
	if err := json.Unmarshal(dataByte, &headDepart); err != nil {
		log.Errorln(pkgName, err, "Json unmarshall head department error")
		return headDepart, nil
	}
	return headDepart, nil
}
