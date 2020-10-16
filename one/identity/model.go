package identity

import uuid "github.com/satori/go.uuid"

//
// Identity
//

// Authentication result model
type AuthenticationResult struct {
	TokenType    string         `json:"token_type"`
	ExpiresIn    int            `json:"expires_in"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	AccountID    string         `json:"account_id"`
	Result       string         `json:"result"`
	Username     string         `json:"username"`
	Profile      AccountProfile `json:"profile"`
}

// Account profile model
type AccountProfile struct {
	ID                 string          `json:"id"`
	FirstNameTH        string          `json:"first_name_th"`
	LastNameTH         string          `json:"last_name_th"`
	FirstNameENG       string          `json:"first_name_eng"`
	LastNameENG        string          `json:"last_name_eng"`
	TitleTH            string          `json:"account_title_th"`
	TitleENG           string          `json:"account_title_eng"`
	IDCardType         string          `json:"id_card_type"`
	IDCardTypeNumber   string          `json:"id_card_num"`
	IDCardHashed       string          `json:"hash_id_card_num"`
	AccountCategory    string          `json:"account_category"`
	AccountSubCategory string          `json:"account_sub_category"`
	ThaiEmail1         string          `json:"thai_email"`
	ThaiEmail2         string          `json:"thai_email2"`
	StatusCD           string          `json:"status_cd"`
	BirthDate          string          `json:"birth_date"`
	StatusDate         string          `json:"status_dt"`
	RegisterDate       string          `json:"register_dt"`
	AddressID          string          `json:"address_id"`
	CreatedAt          string          `json:"created_at"`
	CreatedBy          string          `json:"created_by"`
	UpdatedAt          string          `json:"updated_at"`
	UpdatedBy          string          `json:"updated_by"`
	Reason             string          `json:"reason"`
	TelephoneNumber    string          `json:"tel_no"`
	NameOnDocTH        string          `json:"name_on_document_th"`
	NameOnDocENG       string          `json:"name_on_document_eng"`
	Mobile             []AccountMobile `json:"mobile"`
	Email              []AccountEmail  `json:"email"`
	Address            []string        `json:"address"`
	AccountAttr        []string        `json:"account_attribute"`
	Status             string          `json:"status"`
	LastUpdate         string          `json:"last_update"`
	Employee           *Employee       `json:"has_employee_detail"`
}

type AccountMobile struct {
	ID           string             `json:"id"`
	MobileNumber string             `json:"mobile_no"`
	CreatedAt    string             `json:"created_at"`
	CreatedBy    string             `json:"created_by"`
	UpdatedAt    string             `json:"updated_at"`
	UpdatedBy    string             `json:"updated_by"`
	DeletedAt    string             `json:"deleted_at"`
	MobilePivot  AccountMobilePivot `json:"pivot"`
}

type AccountMobilePivot struct {
	AccountID   string `json:"account_id"`
	MobileID    string `json:"mobile_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	StatusCD    string `json:"status_cd"`
	PrimaryFlag string `json:"primary_flg"`
	ConfirmFlag string `json:"mobile_confirm_flg"`
	ConfirmDate string `json:"mobile_confirm_dt"`
}

type AccountEmail struct {
	ID         string            `json:"id"`
	Email      string            `json:"email"`
	CreatedAt  string            `json:"created_at"`
	CreatedBy  string            `json:"created_by"`
	UpdatedAt  string            `json:"updated_at"`
	UpdatedBy  string            `json:"updated_by"`
	DeletedBy  string            `json:"deleted_at"`
	EmailPivot AccountEmailPivot `json:"pivot"`
}

type AccountEmailPivot struct {
	AccountID   string `json:"account_id"`
	EmailID     string `json:"email_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	StatusCD    string `json:"status_cd"`
	PrimaryFlag string `json:"primary_flg"`
	ConfirmFlag string `json:"email_confirm_flg"`
	ConfirmDate string `json:"email_confirm_dt"`
}

//
// Organize Chart
//

type OrgClient struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"type"`
	ApiEndpoint  string `json:"api_endpoint"`
}

type OrgApiResult struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"errorMessage"`
	Code   int         `json:"code"`
}

type Department struct {
	Id           uuid.UUID  `json:"id"`
	Name         string     `json:"dept_name"`
	ParentDeptId *uuid.UUID `json:"parent_dept_id"`
	Accounts     *Employee  `json:"has_account"`
}

type TeamMember struct {
	Id           uuid.UUID   `json:"id"`
	Name         string      `json:"dept_name"`
	ParentDeptId *uuid.UUID  `json:"parent_dept_id"`
	Accounts     *[]Employee `json:"has_account"`
}

type HeadDepartment struct {
	Id           uuid.UUID   `json:"id"`
	Name         string      `json:"dept_name"`
	ParentDeptId *uuid.UUID  `json:"parent_dept_id"`
	Accounts     *[]Employee `json:"has_account"`
}

type Employee struct {
	Id         uuid.UUID       `json:"id"`
	AccountId  string          `json:"account_id"`
	BizId      string          `json:"biz_id"`
	Email      string          `json:"email"`
	EmployeeId string          `json:"employee_id"`
	Account    *AccountProfile `json:"account"`
	Employee   *Employee       `json:"employee"`
	Position   string          `json:"position"`
	PositionId uuid.UUID       `json:"role_id"`
}
