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

func NewChatBot(botId string, token string, tokenType string, apiEndpoint *string) Chat {
	ep := ProductionEndpoint
	if apiEndpoint != nil && *apiEndpoint != "" {
		ep = *apiEndpoint
	}
	return Chat{
		BotId:       botId,
		Token:       token,
		TokenType:   tokenType,
		ApiEndpoint: ep,
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
	r, err := c.send(http.MethodPost, c.url("/manage/api/v1/getprofile"), body)
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
	r, err := c.send(http.MethodPost, c.url("/message/api/v1/searchfriend"), body)
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

func (c *Chat) GetOneChatAccessToken(oneChatToken string) (string, error) {
	msg := struct {
		OneChatToken string `json:"onechat_token"`
	}{
		OneChatToken: oneChatToken,
	}
	body, _ := json.Marshal(&msg)
	r, err := c.sendWithCustomToken(http.MethodPost, c.url("/event/api/v1/accesstoken_by_onechattoken"), body, "")
	if err != nil {
		return "", err
	}
	result := struct {
		AccessToken string `json:"access_token"`
		Status      string `json:"status"`
	}{}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
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
	_, err := c.send(http.MethodPost, c.url("/message/api/v1/push_message"), body)
	return err
}

func (c *Chat) PushWebView(to string, title string, detail string, img string, choices []Choice, customNotify *string) error {
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
				Image:   img,
				Title:   title,
				Detail:  detail,
				Choices: choices,
			},
		},
	}

	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/message/api/v1/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Chat) PushLinkTemplate(to string, title string, detail string, img string, choices []Choice, customNotify *string) error {
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
				Image:   img,
				Title:   title,
				Detail:  detail,
				Choices: choices,
			},
		},
	}

	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/message/api/v1/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Chat) GetSharedToken(oneChatToken string) (string, error) {
	msg := struct {
		OneChatToken string `json:"onechat_token"`
	}{
		OneChatToken: oneChatToken,
	}
	body, _ := json.Marshal(&msg)
	r, err := c.send(http.MethodPost, c.url("/event/api/v1/sharedtoken_by_onechattoken"), body)
	if err != nil {
		return "", err
	}
	result := struct {
		Data struct {
			SharedToken string `json:"shared_token"`
		} `json:"data"`
		Status string `json:"status"`
	}{}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return "", err
	}
	return result.Data.SharedToken, nil
}

func (c *Chat) CloseWebView(to string) error {
	pushBody := struct {
		UserId string `json:"user_id"`
		BotId  string `json:"bot_id"`
	}{
		UserId: to,
		BotId:  c.BotId,
	}
	body, _ := json.Marshal(&pushBody)
	r, err := c.send(http.MethodPost, c.url("/message/api/v2/disable_webview"), body)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Chat) send(method string, url string, body []byte) (requests.Response, error) {
	return c.sendWithCustomToken(method, url, body, c.Token)
}

func (c *Chat) sendWithCustomToken(method string, url string, body []byte, token string) (requests.Response, error) {
	headers := map[string]string{
		httputil.HeaderContentType: "application/json",
	}
	if token != "" {
		headers[httputil.HeaderAuthorization] = fmt.Sprintf("%s %s", c.TokenType, token)
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
