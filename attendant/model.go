package attendant

// EmployeeDetail struct
type EmployeeDetail struct {
	Company    string `json:"company"`
	Department string `json:"department"`
	Email      string `json:"email"`
	EmployeeID string `json:"employee_id"`
	Position   string `json:"position"`
	TaxID      string `json:"tax_id"`
	Telephone  string `json:"telephone"`
}

// AccountProfile struct
type AccountProfile struct {
	AccountID      string           `json:"account_id"`
	EmployeeDetail []EmployeeDetail `json:"employee_detail"`
	FirstNameTH    string           `json:"first_name_th"`
	LastNameTH     string           `json:"last_name_th"`
	NicknameTH     string           `json:"nick_name_th"`
	OneEmail       string           `json:"one_email"`
	TelNO          string           `json:"tel_no"`
	TitleTH        string           `json:"title_th"`
}

// HasAccount struct
type HasAccount struct {
	AccountID      string `json:"account_id"`
	AccountTitleEN string `json:"account_title_en"`
	EmployeeID     string `json:"employee_id"`
	FirstNameEN    string `json:"first_name_en"`
	LastNameEN     string `json:"last_name_en"`
	Position       string `json:"position"`
}

// Departments struct
type Departments struct {
	DeptName   string       `json:"dept_name"`
	HasAccount []HasAccount `json:"has_account"`
	UUID       string       `json:"uid"`
}

// CompanyDetail struct
type CompanyDetail struct {
	CompanyName string        `json:"company_name"`
	Departments []Departments `json:"departments"`
	TaxNo       string        `json:"tax_no"`
}

// Companies struct
type Companies struct {
	Name  string `json:"name"`
	TaxNo string `json:"tax_no"`
}
