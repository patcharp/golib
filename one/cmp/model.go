package cmp

const (
	ProductionEndpoint = "https://apigw-86xkfxhy.appsdigit.inet.co.th"
)

type Client struct {
	ApiEndpoint string
}

type ConsentRequestForm struct {
	TaxId         string `json:"tax_id"`
	ShareToken    string `json:"shared_token"`
	ServiceId     string `json:"service_id"`
	ConsentFormId string `json:"consent_form_id"`
}

type StatusResult struct {
	Result       string     `json:"result"`
	Data         StatusInfo `json:"data"`
	ErrorMessage string     `json:"errorMessage"`
	Code         int        `json:"code"`
}

type StatusInfo struct {
	Status   string       `json:"status"`
	Response ChatResponse `json:"response"`
}

type ChatResponse struct {
	To           string   `json:"to"`
	CustomNotify string   `json:"custom_notification"`
	Elements     Elements `json:"elements"`
}

type Elements struct {
	Image   string `json:"image"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
	Choices Choice `json:"choice"`
}

type Choice struct {
	Label        string `json:"label"`
	Type         string `json:"type"`
	Url          string `json:"url"`
	Size         string `json:"size"`
	Sign         string `json:"sign"`
	OneChatToken bool   `json:"onechat_token"`
}
