package identity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patcharp/golib/requests"
	"net/http"
)

const (
	IdProductionEndpoint = "https://one.th"

	grantTypePassword     = "password"
	grantTypeCode         = "authorization_code"
	grantTypeRefreshToken = "refresh_token"

	MaxTimeOut = 10
)

type IdCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Endpoint string `json:"endpoint"`
}

type Identity struct {
	ApiEndpoint  string
	ClientId     string
	ClientSecret string
}

func NewIdentity(clientId string, clientSecret string) Identity {
	return Identity{
		ApiEndpoint:  IdProductionEndpoint,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

func (id *Identity) SetEndpoint(ep string) {
	id.ApiEndpoint = ep
}

func (id *Identity) Login(username string, password string, profile bool) (AuthenticationResult, error) {
	var result AuthenticationResult
	body, _ := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}{
		ClientID:     id.ClientId,
		ClientSecret: id.ClientSecret,
		GrantType:    grantTypePassword,
		Username:     username,
		Password:     password,
	})
	headers := map[string]string{
		echo.HeaderContentType: "application/json",
	}
	r, err := requests.Post(id.url("/api/oauth/getpwd"), headers, bytes.NewBuffer(body), MaxTimeOut)
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
		result.Profile, err = id.profile(result.TokenType, result.AccessToken)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (id *Identity) RefreshNewToken(refreshToken string) (AuthenticationResult, error) {
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
		ClientID:     id.ClientId,
		ClientSecret: id.ClientSecret,
		GrantType:    grantTypeRefreshToken,
		RefreshToken: refreshToken,
	})
	headers := map[string]string{
		echo.HeaderContentType: "application/json",
	}
	resp, err := requests.Post(id.url("/api/oauth/get_refresh_token"), headers, bytes.NewBuffer(body), MaxTimeOut)
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

func (id *Identity) VerifyAuthorizedCode(authCode string) (AuthenticationResult, error) {
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
	}{
		GrantType:    grantTypeCode,
		ClientID:     id.ClientId,
		ClientSecret: id.ClientSecret,
		Code:         authCode,
		Scope:        "",
	})
	headers := map[string]string{
		echo.HeaderContentType: "application/json",
	}
	if err != nil {
		return result, err
	}
	resp, err := requests.Post(id.url("/oauth/token"), headers, bytes.NewBuffer(reqJson), MaxTimeOut)
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

func (id *Identity) profile(tokenType string, accessToken string) (AccountProfile, error) {
	var profile AccountProfile
	if tokenType == "" || accessToken == "" {
		return profile, errors.New("login required")
	}
	headers := map[string]string{
		echo.HeaderAuthorization: fmt.Sprintf("%s %s", tokenType, accessToken),
	}
	rawResponse, err := requests.Get(id.url("/api/account"), headers, nil, MaxTimeOut)
	if err != nil {
		return profile, err
	}
	return profile, json.Unmarshal(rawResponse.Body, &profile)
}

func (id *Identity) url(path string) string {
	return fmt.Sprintf("%s%s", id.ApiEndpoint, path)
}

func (id *Identity) GetLoginUrl() string {
	return id.url(fmt.Sprintf("/api/oauth/getcode?client_id=%s&response_type=%s&scope=%s", id.ClientId, "code", ""))
}

func (id *Identity) RedirectToLoginUrl(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, id.GetLoginUrl(), http.StatusFound)
}
