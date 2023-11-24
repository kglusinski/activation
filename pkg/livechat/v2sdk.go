package livechat

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"activation_flow/internal/security"
)

const apiURL = "https://accounts.livechat.com"

type Client struct {
	client *http.Client
	cfg    Config
}

func NewClient(cfg Config) *Client {
	return &Client{
		client: &http.Client{},
		cfg:    cfg,
	}
}

// GetToken sends request to https://accounts.livechat.com/v2/token
func (c *Client) GetToken(token string) (*security.Token, error) {
	resource := "/v2/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", c.cfg.ClientID)
	data.Set("code", token)
	data.Set("client_secret", c.cfg.ClientSecret)
	data.Set("redirect_uri", c.cfg.RedirectUrl)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	log.Println("data: ", data.Encode())

	log.Println("sending post to: ", urlStr)
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(r)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	log.Println("success response: ", resp.Status)

	var tokenResponse security.Token
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

// GetUser sends request to https://accounts.livechat.com/v2/accounts/me
func (c *Client) GetUser(accessToken security.Token) (*security.User, error) {
	resource := "/v2/accounts/me"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	r, _ := http.NewRequest(http.MethodGet, urlStr, nil)
	r.Header.Add("Authorization", "Bearer "+accessToken.AccessToken)
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(r)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var userResponse security.User
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse, nil
}
