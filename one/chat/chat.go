package chat

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/patcharp/golib/v2/requests"
	"github.com/patcharp/golib/v2/util/httputil"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	r, err := c.send(http.MethodPost, c.url("/go_api/api/v1/get-profile"), body)
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

func (c *Chat) GetFriendList() ([]Friend, error) {
	body, _ := json.Marshal(&struct {
		BotId string `json:"bot_id"`
	}{
		BotId: c.BotId,
	})
	r, err := c.send(http.MethodPost, c.url("/manage/api/v1/getlistroom"), body)
	if err != nil {
		return nil, err
	}
	if r.Code != http.StatusOK {
		return nil, errors.New(fmt.Sprint("OneChat return code", r.Code, string(r.Body)))
	}
	var fiendList FriendList
	if err := json.Unmarshal(r.Body, &fiendList); err != nil {
		return nil, err
	}
	return fiendList.ListFriend, nil
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

func (c *Chat) PushBroadcastMessage(to []string, msg string) error {
	if len(to) > 100 {
		return errors.New("receiver should not more than 100 accounts")
	}
	pushMessage := struct {
		To      []string `json:"to"`
		BotId   string   `json:"bot_id"`
		Message string   `json:"message"`
	}{
		To:      to,
		BotId:   c.BotId,
		Message: msg,
	}
	body, _ := json.Marshal(&pushMessage)
	_, err := c.send(http.MethodPost, c.url("/bc_msg/api/v1/broadcast_group"), body)
	return err
}
func (c *Chat) PushBroadcastFromByte(to []string, file []byte) error {
	fileBuff := bytes.Buffer{}
	fileBuff.Write(file)
	hash := sha256.New()
	hash.Write(file)
	kind, _ := filetype.Match(file)
	filename := fmt.Sprintf("%s.%s", base64.URLEncoding.EncodeToString(hash.Sum(nil)), kind.Extension)
	if err := ioutil.WriteFile(filename, file, 0644); err != nil {
		return err
	}
	defer os.Remove(filename)
	return c.PushBroadcastFromFile(to, filename)
}

func (c *Chat) PushBroadcastFromFile(to []string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("bot_id", c.BotId)
	for _, account := range to {
		_ = writer.WriteField("to", account)
	}
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return err
	}
	headers := map[string]string{
		httputil.HeaderAuthorization: fmt.Sprintf("%s %s", c.TokenType, c.Token),
		httputil.HeaderContentType:   writer.FormDataContentType(),
	}
	r, err := c.sendWithCustomHeader(http.MethodPost, c.url("/bc_msg/api/v1/broadcast_group_file"), body, headers)
	if err != nil {
		return err
	}
	if r.Code != http.StatusOK {
		return errors.New(fmt.Sprint("server return error status", string(r.Body)))
	}
	return nil
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

func (c *Chat) PushQuickReplyTextType(to string, message string, quickReplies []QuickReplyTextType) error {
	body, _ := json.Marshal(&struct {
		To         string               `json:"to"`
		BotId      string               `json:"bot_id"`
		Message    string               `json:"message"`
		QuickReply []QuickReplyTextType `json:"quick_reply"`
	}{
		To:         to,
		BotId:      c.BotId,
		Message:    message,
		QuickReply: quickReplies,
	})
	resp, err := c.send(http.MethodPost, c.url("/message/api/v1/push_quickreply"), body)
	if err != nil {
		return err
	}
	if resp.Code != http.StatusOK {
		return errors.New(fmt.Sprintln("OneChat server return not ok code", resp.Code, string(resp.Body)))
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
	headers := map[string]string{
		httputil.HeaderContentType:   "application/json",
		httputil.HeaderAuthorization: fmt.Sprintf("%s %s", c.TokenType, c.Token),
	}
	return c.sendWithCustomHeader(method, url, bytes.NewBuffer(body), headers)
}

func (c *Chat) sendWithCustomToken(method string, url string, body []byte, token string) (requests.Response, error) {
	headers := map[string]string{
		httputil.HeaderContentType:   "application/json",
		httputil.HeaderAuthorization: fmt.Sprintf("%s %s", c.TokenType, token),
	}
	return c.sendWithCustomHeader(method, url, bytes.NewBuffer(body), headers)
}

func (c *Chat) sendWithCustomHeader(method string, url string, body io.Reader, headers map[string]string) (requests.Response, error) {
	r, err := requests.Request(method, url, headers, body, 0)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (c *Chat) url(path string) string {
	return fmt.Sprintf("%s%s", c.ApiEndpoint, path)
}
