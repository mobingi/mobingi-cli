package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/pkg/errors"
)

const (
	BASE_API_URL      = "https://api.mobingi.com"
	BASE_REGISTRY_URL = "https://registry.mobingi.com"
	SESHA3_URL        = "https://sesha3.labs.mobingi.com"
)

type authPayload struct {
	ClientId     string      `json:"client_id,omitempty"`
	ClientSecret string      `json:"client_secret,omitempty"`
	GrantType    string      `json:"grant_type,omitempty"`
	Username     interface{} `json:"username,omitempty"`
	Password     interface{} `json:"password,omitempty"`
}

type Config struct {
	// ClientId is your Mobingi client id. If empty, it will look for
	// MOBINGI_CLIENT_ID environment variable.
	ClientId string

	// ClientSecret is your Mobingi client secret. If empty, it will look for
	// MOBINGI_CLIENT_SECRET environment variable.
	ClientSecret string

	// Username is your Mobingi subuser name. If empty, it means the login grant
	// type is 'client_credentials'.
	Username string

	// Password is your Mobingi subuser password. Cannot be empty when Username
	// is not empty.
	Password string

	// AccessToken is your API access token. By default, session will get an
	// access token based on ClientId and ClientSecret. If this is set however,
	// session will use this token instead.
	AccessToken string

	// ApiVersion is the API version to be used in the session where this config
	// is associated with. If zero, it will default to the latest version.
	ApiVersion int

	// BaseApiUrl is the base API URL for this session. Default is the latest
	// production endpoint.
	BaseApiUrl string

	// BaseRegistryUrl is the base URL for Mobingi Docker Registry. Default is the
	// latest production endpoint.
	BaseRegistryUrl string

	// Sesha3Url is the base URL for sesha3. Default is the latest production endpoint.
	Sesha3Url string

	// HttpClientConfig will set the config for the session's http client. Do not
	// set if you want to use http client defaults.
	HttpClientConfig *client.Config
}

type Session struct {
	Config      *Config
	AccessToken string
}

func (s *Session) ApiEndpoint() string {
	return fmt.Sprintf("%s/v%d", s.Config.BaseApiUrl, s.Config.ApiVersion)
}

func (s *Session) RegistryEndpoint() string {
	return fmt.Sprintf("%s/v2", s.Config.BaseRegistryUrl)
}

func (s *Session) Sesha3Endpoint() string {
	return s.Config.Sesha3Url
}

func (s *Session) SimpleAuthRequest(m, u string, body io.Reader) *http.Request {
	req, err := http.NewRequest(m, u, body)
	if err != nil {
		return nil
	}

	req.Header.Add("Authorization", "Bearer "+s.AccessToken)
	return req
}

func (s *Session) getAccessToken() (string, error) {
	var token string
	var p authPayload

	if s.Config.Username != "" {
		if s.Config.Password == "" {
			return token, errors.New("password cannot be empty")
		}

		p = authPayload{
			ClientId:     s.Config.ClientId,
			ClientSecret: s.Config.ClientSecret,
			GrantType:    "password",
			Username:     s.Config.Username,
			Password:     s.Config.Password,
		}
	} else {
		p = authPayload{
			ClientId:     s.Config.ClientId,
			ClientSecret: s.Config.ClientSecret,
			GrantType:    "client_credentials",
		}
	}

	payload, err := json.Marshal(p)
	r, err := http.NewRequest(
		http.MethodPost,
		s.ApiEndpoint()+"/access_token",
		bytes.NewBuffer(payload))
	if err != nil {
		return token, errors.Wrap(err, "new request failed")
	}

	r.Header.Add("Content-Type", "application/json")

	var c client.HttpClient

	if s.Config.HttpClientConfig != nil {
		c = client.NewSimpleHttpClient(s.Config.HttpClientConfig)
	} else {
		c = client.NewSimpleHttpClient()
	}

	resp, body, err := c.Do(r)
	if err != nil {
		return token, errors.Wrap(err, "do failed")
	}

	if (resp.StatusCode / 100) != 2 {
		return token, errors.New(resp.Status)
	}

	var m map[string]interface{}
	if err = json.Unmarshal(body, &m); err != nil {
		return token, errors.Wrap(err, "unmarshal failed")
	}

	t, found := m["access_token"]
	if !found {
		return token, fmt.Errorf("cannot find access token")
	}

	token = fmt.Sprintf("%s", t)
	return token, nil
}

func New(cnf ...*Config) (*Session, error) {
	c := &Config{
		ClientId:        os.Getenv("MOBINGI_CLIENT_ID"),
		ClientSecret:    os.Getenv("MOBINGI_CLIENT_SECRET"),
		Username:        os.Getenv("MOBINGI_USERNAME"),
		Password:        os.Getenv("MOBINGI_PASSWORD"),
		ApiVersion:      3,
		BaseApiUrl:      BASE_API_URL,
		BaseRegistryUrl: BASE_REGISTRY_URL,
		Sesha3Url:       SESHA3_URL,
	}

	if len(cnf) > 0 {
		if cnf[0] != nil {
			if cnf[0].ClientId != "" {
				c.ClientId = cnf[0].ClientId
			}

			if cnf[0].ClientSecret != "" {
				c.ClientSecret = cnf[0].ClientSecret
			}

			if cnf[0].AccessToken != "" {
				c.AccessToken = cnf[0].AccessToken
			}

			if cnf[0].Username != "" {
				c.Username = cnf[0].Username
			}

			if cnf[0].Password != "" {
				c.Password = cnf[0].Password
			}

			if cnf[0].ApiVersion > 0 {
				c.ApiVersion = cnf[0].ApiVersion
			}

			if cnf[0].BaseApiUrl != "" {
				c.BaseApiUrl = cnf[0].BaseApiUrl
			}

			if cnf[0].BaseRegistryUrl != "" {
				c.BaseRegistryUrl = cnf[0].BaseRegistryUrl
			}

			if cnf[0].Sesha3Url != "" {
				c.Sesha3Url = cnf[0].Sesha3Url
			}

			if cnf[0].HttpClientConfig != nil {
				c.HttpClientConfig = cnf[0].HttpClientConfig
			}
		}
	}

	s := &Session{Config: c}
	if c.AccessToken != "" {
		s.AccessToken = c.AccessToken
	} else {
		token, err := s.getAccessToken()
		if err != nil {
			return s, errors.Wrap(err, "get access token failed")
		}

		s.AccessToken = token
	}

	return s, nil
}
