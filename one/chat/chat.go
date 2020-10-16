package chat

const ProductionEndpoint = "https://chat-api.one.th/message/api/v1"

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
	ProfilePicture string `json:"profilepicture"`
}

func NewChatBot(botId string, token string, tokenType string) Chat {
	return Chat{
		BotId:       botId,
		Token:       token,
		TokenType:   tokenType,
		ApiEndpoint: ProductionEndpoint,
	}
}
