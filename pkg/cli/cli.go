package cli

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

const (
	BaseUrl = "https://apidev.mobingi.com"
)

type Config struct {
	RootUrl   string                // base url + version
	Requester *gorequest.SuperAgent // http requester
}

func New(ver string) *Config {
	return &Config{
		RootUrl:   BaseUrl + "/" + ver,
		Requester: gorequest.New(),
	}
}

func (c *Config) GetSafe(url, token string) (gorequest.Response, []byte, []error) {
	return c.Requester.Get(url).Set("Authorization", "Bearer "+token).EndBytes()
}

func (c *Config) PostJSON(url string, i interface{}) (gorequest.Response, []byte, []error) {
	e := make([]error, 0)
	payload, err := json.Marshal(i)
	if err != nil {
		e = append(e, err)
		return nil, nil, e
	}

	return c.Requester.Post(url).Send(string(payload)).EndBytes()
}

func (c *Config) DeleteSafe(url, token string) (gorequest.Response, []byte, []error) {
	return c.Requester.Delete(url).Set("Authorization", "Bearer "+token).EndBytes()
}
