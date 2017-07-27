package registry

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/credentials"
)

type TokenCredentials struct {
	UserPass *credentials.UserPass
	Service  string
	Scope    string
}

type TokenParams struct {
	Base       string
	ApiVersion string
	TokenCreds *TokenCredentials
}

func GetRegistryToken(tp *TokenParams) ([]byte, string, error) {
	var token string
	if tp.TokenCreds == nil {
		return nil, token, fmt.Errorf("credentials cannot be nil")
	}

	if tp.TokenCreds.UserPass == nil {
		return nil, token, fmt.Errorf("credentials cannot be nil")
	}

	_, err := tp.TokenCreds.UserPass.EnsureInput(false)
	if err != nil {
		return nil, token, err
	}

	c := client.NewClient(&client.Config{
		RootUrl:    tp.Base,
		ApiVersion: tp.ApiVersion,
	})

	v := url.Values{}
	v.Add("service", tp.TokenCreds.Service)
	v.Add("scope", tp.TokenCreds.Scope)

	body, err := c.BasicAuthGet(
		"/docker/token",
		tp.TokenCreds.UserPass.Username,
		tp.TokenCreds.UserPass.Password,
		&v,
	)

	if err != nil {
		return nil, token, err
	}

	if len(body) <= 0 {
		return nil, token, fmt.Errorf("empty return")
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, token, err
	}

	t, found := m["token"]
	if !found {
		return nil, token, fmt.Errorf("cannot find token")
	}

	token = fmt.Sprintf("%s", t)
	return body, token, nil
}
