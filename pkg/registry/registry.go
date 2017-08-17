package registry

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/mobingi/mobingi-cli/client"
	"github.com/mobingi/mobingi-cli/pkg/credentials"
	"github.com/pkg/errors"
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
		return nil, token, errors.Wrap(err, "ensure input failed")
	}

	c := client.NewClient(&client.Config{
		RootUrl:    tp.Base,
		ApiVersion: tp.ApiVersion,
	})

	values := url.Values{}
	values.Add("service", tp.TokenCreds.Service)
	values.Add("scope", tp.TokenCreds.Scope)

	body, err := c.BasicAuthGet(
		"/docker/token",
		tp.TokenCreds.UserPass.Username,
		tp.TokenCreds.UserPass.Password,
		&values,
	)

	if err != nil {
		return nil, token, errors.Wrap(err, "basic auth get failed")
	}

	if len(body) <= 0 {
		return nil, token, fmt.Errorf("empty return")
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, token, errors.Wrap(err, "unmarshal failed")
	}

	t, found := m["token"]
	if !found {
		return nil, token, fmt.Errorf("cannot find token")
	}

	token = fmt.Sprintf("%s", t)
	return body, token, nil
}
