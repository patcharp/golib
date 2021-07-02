package identity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patcharp/golib/requests"
	"github.com/patcharp/golib/util/httputil"
	"net/http"
)

func NewIdentity(clientId string, clientSecret string, redirectUrl string, apiEndpoint *string) Client {
	ep := ProductionEndpoint
	if apiEndpoint != nil && *apiEndpoint != "" {
		ep = *apiEndpoint
	}
	return Client{
		ApiEndpoint:  ep,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectUrl:  redirectUrl,
	}
}

func (c *Client) Login(username string, password string, profile bool) (AuthenticationResult, error) {
	var result AuthenticationResult
	body, _ := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		GrantType:    grantTypePassword,
		Username:     username,
		Password:     password,
	})
	headers := map[string]string{
		httputil.HeaderContentType: httputil.MIMEApplicationJSON,
	}
	r, err := requests.Post(c.url("/api/oauth/getpwd"), headers, bytes.NewBuffer(body), MaxTimeOut)
	if err != nil {
		return result, err
	}
	if r.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", r.Code, string(r.Body)))
	}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return result, err
	}
	if profile {
		result.Profile, err = c.GetProfile(result.TokenType, result.AccessToken)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (c *Client) RefreshNewToken(refreshToken string) (AuthenticationResult, error) {
	var result AuthenticationResult
	if refreshToken == "" {
		return result, errors.New("unauthorized identity")
	}
	body, _ := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
	}{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		GrantType:    grantTypeRefreshToken,
		RefreshToken: refreshToken,
	})
	headers := map[string]string{
		httputil.HeaderContentType: httputil.MIMEApplicationJSON,
	}
	resp, err := requests.Post(c.url("/api/oauth/get_refresh_token"), headers, bytes.NewBuffer(body), MaxTimeOut)
	if err != nil {
		return result, err
	}
	if resp.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", resp.Code, string(resp.Body)))
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (c *Client) VerifyAuthorizedCode(authCode string, profile bool) (AuthenticationResult, error) {
	var result AuthenticationResult
	if authCode == "" {
		return result, errors.New("authorization code required")
	}
	reqJson, err := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		Scope        string `json:"scope"`
		RedirectUri  string `json:"redirect_uri"`
	}{
		GrantType:    grantTypeCode,
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Code:         authCode,
		Scope:        "",
		RedirectUri:  c.RedirectUrl,
	})
	headers := map[string]string{
		httputil.HeaderContentType: httputil.MIMEApplicationJSON,
	}
	if err != nil {
		return result, err
	}
	resp, err := requests.Post(c.url("/oauth/token"), headers, bytes.NewBuffer(reqJson), MaxTimeOut)
	if err != nil {
		return result, err
	}
	if resp.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", resp.Code, string(resp.Body)))
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return result, err
	}
	if profile {
		result.Profile, err = c.GetProfile(result.TokenType, result.AccessToken)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (c *Client) LoginWithScope(username string, password string, profile bool, scope string) (AuthenticationResult, error) {
	var result AuthenticationResult
	body, _ := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		Scope        string `json:"scope"`
	}{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		GrantType:    grantTypePassword,
		Username:     username,
		Password:     password,
		Scope:        scope,
	})
	headers := map[string]string{
		httputil.HeaderContentType: httputil.MIMEApplicationJSON,
	}
	r, err := requests.Post(c.url("/api/oauth/getpwd"), headers, bytes.NewBuffer(body), MaxTimeOut)
	if err != nil {
		return result, err
	}
	if r.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", r.Code, string(r.Body)))
	}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return result, err
	}
	if profile {
		result.Profile, err = c.GetProfile(result.TokenType, result.AccessToken)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (c *Client) GetProfile(tokenType string, accessToken string) (AccountProfile, error) {
	var profile AccountProfile
	if tokenType == "" || accessToken == "" {
		return profile, errors.New("login required")
	}
	headers := map[string]string{
		httputil.HeaderAuthorization: fmt.Sprintf("%s %s", tokenType, accessToken),
	}
	rawResponse, err := requests.Get(c.url("/api/account"), headers, nil, MaxTimeOut)
	if err != nil {
		return profile, err
	}
	return profile, json.Unmarshal(rawResponse.Body, &profile)
}

func (c *Client) url(path string) string {
	return fmt.Sprintf("%s%s", c.ApiEndpoint, path)
}

func (c *Client) GetLoginUrl() string {
	return c.url(fmt.Sprintf(
		"/api/oauth/getcode?client_id=%s&response_type=code&scope=&redirect_uri=%s",
		c.ClientId,
		c.RedirectUrl,
	))
}

func (c *Client) RedirectToLoginUrl(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, c.GetLoginUrl(), http.StatusFound)
}
