package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/mobingilabs/mocli/client"
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

func GetRegistryToken(tp *TokenParams) ([]byte, string, error) {
	if tp.TokenCreds == nil {
		return nil, "", fmt.Errorf("credentials cannot be nil")
	}

	if tp.TokenCreds.UserPass == nil {
		return nil, "", fmt.Errorf("credentials cannot be nil")
	}

	_, err := tp.TokenCreds.UserPass.EnsureInput(false)
	if err != nil {
		return nil, "", err
	}

	var u *url.URL
	u, err = url.Parse(tp.Base)
	if err != nil {
		return nil, "", err
	}

	u.Path += "/" + tp.ApiVersion + "/docker/token"
	v := url.Values{}
	v.Add("service", tp.TokenCreds.Service)
	v.Add("scope", tp.TokenCreds.Scope)
	u.RawQuery = v.Encode()

	// TODO: use our client library
	c := &http.Client{
		Timeout: time.Second * time.Duration(client.Timeout),
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, "", err
	}

	req.SetBasicAuth(tp.TokenCreds.UserPass.Username, tp.TokenCreds.UserPass.Password)
	if d.Verbose {
		d.Info(fmt.Sprintf("Get token for subuser '%s' with service '%s' and scope '%s'.",
			tp.TokenCreds.UserPass.Username, tp.TokenCreds.Service, tp.TokenCreds.Scope))
	}

	if d.Verbose {
		for n, h := range req.Header {
			d.Info(fmt.Sprintf("[GETTOKEN-REQUEST] %s: %s", n, h))
		}
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, "", err
	}

	if d.Verbose {
		for n, h := range resp.Header {
			d.Info(fmt.Sprintf("[GETTOKEN-RESPONSE] %s: %s", n, h))
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if len(body) <= 0 {
		return nil, "", fmt.Errorf("empty return")
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
