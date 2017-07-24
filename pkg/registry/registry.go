package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
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

func GetRegistryToken(c *TokenParams, verbose bool) ([]byte, string, error) {
	if c.TokenCreds == nil {
		return nil, "", fmt.Errorf("credentials cannot be nil")
	}

	if c.TokenCreds.UserPass == nil {
		return nil, "", fmt.Errorf("credentials cannot be nil")
	}

	_, err := c.TokenCreds.UserPass.EnsureInput(false)
	if err != nil {
		return nil, "", err
	}

	var u *url.URL
	u, err = url.Parse(c.Base)
	if err != nil {
		return nil, "", err
	}

	u.Path += "/" + c.ApiVersion + "/docker/token"
	v := url.Values{}
	v.Add("service", c.TokenCreds.Service)
	v.Add("scope", c.TokenCreds.Scope)
	u.RawQuery = v.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, "", err
	}

	req.SetBasicAuth(c.TokenCreds.UserPass.Username, c.TokenCreds.UserPass.Password)
	if d.Verbose {
		d.Info(fmt.Sprintf("Get token for subuser '%s' with service '%s' and scope '%s'.",
			c.TokenCreds.UserPass.Username, c.TokenCreds.Service, c.TokenCreds.Scope))
	}

	if verbose {
		for n, h := range req.Header {
			d.Info(fmt.Sprintf("[in] %s: %s", n, h))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	if verbose {
		for n, h := range resp.Header {
			d.Info(fmt.Sprintf("[out] %s: %s", n, h))
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, "", err
	}

	t, found := m["token"]
	if !found {
		return nil, "", fmt.Errorf("cannot find token")
	}

	return body, fmt.Sprintf("%s", t), nil
}
