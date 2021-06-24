package chat

const (
	ProductionEndpoint = "https://chat-api.one.th"

	TemplateTypeWebView = "webview"
	TemplateTypeLink    = "link"
)

type Chat struct {
	BotId       string
	Token       string
	TokenType   string
	ApiEndpoint string
}

type Friend struct {
	OneEmail    string `json:"one_email"`
	UserId      string `json:"user_id"`
	AccountId   string `json:"one_id"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
}

type Profile struct {
	Email          string `json:"email"`
	Nickname       string `json:"nickname"`
	AccountId      string `json:"one_id"`
	Telephone      string `json:"telephone"`
	ProfilePicture string `json:"profilepicture"`
}

type LoginOneChatAccount struct {
	Profile struct {
		AccessToken     string `json:"access_token"`
		AccountTitleEng string `json:"account_title_eng"`
		AccountTitleTh  string `json:"account_title_th"`
		Email           string `json:"email"`
		FirstNameEng    string `json:"first_name_eng"`
		FirstNameTh     string `json:"first_name_th"`
		IdCardNum       string `json:"id_card_num"`
		IdCardType      string `json:"id_card_type"`
		Ip              string `json:"ip"`
		LastNameEng     string `json:"last_name_eng"`
		LastNameTh      string `json:"last_name_th"`
		Loa             string `json:"loa"`
		Name            string `json:"name"`
		NickName        string `jaon:"nickname"`
		OneId           string `json:"one_id"`
		Phone           string `json:"phone"`
		Status          string `json:"status"`
		TokenService    string `json:"tokenservice"`
		TokenUser       string `json:"tokenuser"`
		Type            string `json:"type"`
		Username        string `json:"username"`
	} `json:"profile"`
}

type LoginOneChatResult struct {
	Account    LoginOneChatAccount `json:"account"`
	Business   []string            `json:"business"`
	Government []string            `json:"goverment"`
	Message    string              `json:"message"`
	Status     string              `json:"status"`
}

type Choice struct {
	Label        string `json:"label"`
	Type         string `json:"type"`
	Url          string `json:"url"`
	Size         string `json:"size"`
	Sign         string `json:"sign"`
	OneChatToken string `json:"onechat_token"`
}
type Elements struct {
	Image   string   `json:"image"`
	Title   string   `json:"title"`
	Detail  string   `json:"detail"`
	Choices []Choice `json:"choice"`
}

type QuickReplyTextType struct {
	Label   string `json:"label"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Payload string `json:"payload"`
}
