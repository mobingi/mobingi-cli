package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	d "github.com/mobingilabs/mocli/pkg/debug"
)

type TokenParams struct {
	Base       string
	ApiVersion string
	Username   string
	Password   string
	Account    string
	Service    string
	Scope      string
}

func GetRegistryToken(c *TokenParams, verbose bool) ([]byte, string, error) {
	var u *url.URL
	u, err := url.Parse(c.Base)
	if err != nil {
		return nil, "", err
	}

	u.Path += "/" + c.ApiVersion + "/docker/token"
	v := url.Values{}
	v.Add("account", c.Account)
	v.Add("service", c.Service)
	v.Add("scope", c.Scope)
	u.RawQuery = v.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, "", err
	}

	req.SetBasicAuth(c.Username, c.Password)
	d.Info(fmt.Sprintf("Get token for subuser '%s' with service '%s' and scope '%s'.",
		c.Account, c.Service, c.Scope))

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
