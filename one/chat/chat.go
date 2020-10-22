package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patcharp/golib/requests"
	"github.com/patcharp/golib/util/httputil"
	"net/http"
)

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

func (c *Chat) GetProfile(oneChatToken string) (Profile, error) {
	msg := struct {
		BotId        string `json:"bot_id"`
		OneChatToken string `json:"source"`
	}{
		BotId:        c.BotId,
		OneChatToken: oneChatToken,
	}
	body, _ := json.Marshal(&msg)
	r, err := c.send(http.MethodPost, "https://chat-api.one.th/manage/api/v1/getprofile", body)
	if err != nil {
		return Profile{}, err
	}
	chatProfileResult := struct {
		Data   Profile `json:"data"`
		Status string  `json:"status"`
	}{}
	if err := json.Unmarshal(r.Body, &chatProfileResult); err != nil {
		return Profile{}, err
	}
	return chatProfileResult.Data, nil
}

func (c *Chat) FindFriend(keyword string) (Friend, error) {
	msg := struct {
		BotId   string `json:"bot_id"`
		Keyword string `json:"key_search"`
	}{
		BotId:   c.BotId,
		Keyword: keyword,
	}
	body, _ := json.Marshal(&msg)
	r, err := c.send(http.MethodPost, c.url("/searchfriend"), body)
	if err != nil {
		return Friend{}, err
	}
	chatFriendResult := struct {
		Status string `json:"status"`
		Friend Friend `json:"friend"`
	}{}
	if err := json.Unmarshal(r.Body, &chatFriendResult); err != nil {
		return Friend{}, err
	}
	return chatFriendResult.Friend, nil
}

func (c *Chat) PushTextMessage(to string, msg string, customNotify *string) error {
	pushMessage := struct {
		To           string `json:"to"`
		BotId        string `json:"bot_id"`
		Type         string `json:"type"`
		Message      string `json:"message"`
		CustomNotify string `json:"custom_notification,omitempty"`
	}{
		To:      to,
		BotId:   c.BotId,
		Type:    "text",
		Message: msg,
	}
	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	_, err := c.send(http.MethodPost, c.url("/push_message"), body)
	return err
}

func (c *Chat) PushWebView(to string, label string, title string, detail string, path string, img string, customNotify *string) error {
	type Choice struct {
		Label        string `json:"label"`
		Type         string `json:"type"`
		Url          string `json:"url"`
		Size         string `json:"size"`
		OneChatToken string `json:"onechat_token"`
	}
	type Elements struct {
		Image   string   `json:"image"`
		Title   string   `json:"title"`
		Detail  string   `json:"detail"`
		Choices []Choice `json:"choice"`
	}
	pushMessage := struct {
		To           string     `json:"to"`
		BotId        string     `json:"bot_id"`
		Type         string     `json:"type"`
		CustomNotify string     `json:"custom_notification,omitempty"`
		Elements     []Elements `json:"elements"`
	}{
		To:    to,
		BotId: c.BotId,
		Type:  "template",
		Elements: []Elements{
			{
				Image:  img,
				Title:  title,
				Detail: detail,
				Choices: []Choice{
					{
						Label:        label,
						Type:         "webview",
						Url:          path,
						Size:         "full",
						OneChatToken: "true",
					},
				},
			},
		},
	}

	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Chat) PushLinkTemplate(to string, label string, title string, detail string, path string, img string, customNotify *string) error {
	type Choice struct {
		Label        string `json:"label"`
		Type         string `json:"type"`
		Url          string `json:"url"`
		OneChatToken string `json:"onechat_token"`
	}
	type Elements struct {
		Image   string   `json:"image"`
		Title   string   `json:"title"`
		Detail  string   `json:"detail"`
		Choices []Choice `json:"choice"`
	}
	pushMessage := struct {
		To           string     `json:"to"`
		BotId        string     `json:"bot_id"`
		Type         string     `json:"type"`
		CustomNotify string     `json:"custom_notification,omitempty"`
		Elements     []Elements `json:"elements"`
	}{
		To:    to,
		BotId: c.BotId,
		Type:  "template",
		Elements: []Elements{
			{
				Image:  img,
				Title:  title,
				Detail: detail,
				Choices: []Choice{
					{
						Label:        label,
						Type:         "link",
						Url:          path,
						OneChatToken: "true",
					},
				},
			},
		},
	}

	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Chat) send(method string, url string, body []byte) (requests.Response, error) {
	headers := map[string]string{
		httputil.HeaderContentType:   "application/json",
		httputil.HeaderAuthorization: fmt.Sprintf("%s %s", c.TokenType, c.Token),
	}
	r, err := requests.Request(method, url, headers, bytes.NewBuffer(body), 0)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (c *Chat) url(path string) string {
	return fmt.Sprintf("%s%s", c.ApiEndpoint, path)
}
