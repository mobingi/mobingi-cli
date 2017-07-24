package client

import (
	"github.com/parnurzeal/gorequest"
)

type Client struct {
	requester *gorequest.SuperAgent
	config    *Config
}

func NewClient(cnf *Config) *Client {
	return &Client{
		requester: gorequest.New(),
		config:    cnf,
	}
}

func (c *Client) Get(path string) (gorequest.Response, []byte, []error) {
	return c.requester.Get(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).EndBytes()
}

func (c *Client) PostU(path, payload string) (gorequest.Response, []byte, []error) {
	return c.requester.Post(c.url() + path).Send(payload).EndBytes()
}

func (c *Client) Put(path, payload string) (gorequest.Response, []byte, []error) {
	return c.requester.Put(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).Send(payload).EndBytes()
}

func (c *Client) Del(path string) (gorequest.Response, []byte, []error) {
	return c.requester.Delete(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).EndBytes()
}

func (c *Client) url() string {
	return c.config.RootUrl + "/" + c.config.ApiVersion
}
