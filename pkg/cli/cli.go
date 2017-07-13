package cli

import (
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

func (c *Config) Post(url, payload string) (gorequest.Response, []byte, []error) {
	return c.Requester.Post(url).Send(payload).EndBytes()
}

func (c *Config) PutSafe(url, token string, payload string) (gorequest.Response, []byte, []error) {
	return c.Requester.Put(url).Set("Authorization", "Bearer "+token).Send(payload).EndBytes()
}

func (c *Config) DeleteSafe(url, token string) (gorequest.Response, []byte, []error) {
	return c.Requester.Delete(url).Set("Authorization", "Bearer "+token).EndBytes()
}
