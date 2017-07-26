package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/parnurzeal/gorequest"
)

type GrClient struct {
	requester *gorequest.SuperAgent
	config    *Config
}

func NewGrClient(cnf *Config) *GrClient {
	return &GrClient{
		requester: gorequest.New(),
		config:    cnf,
	}
}

func (c *GrClient) Get(path string) (gorequest.Response, []byte, []error) {
	return c.requester.Get(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).EndBytes()
}

func (c *GrClient) PostU(path, payload string) (gorequest.Response, []byte, []error) {
	return c.requester.Post(c.url() + path).Send(payload).EndBytes()
}

func (c *GrClient) Put(path, payload string) (gorequest.Response, []byte, []error) {
	return c.requester.Put(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).Send(payload).EndBytes()
}

func (c *GrClient) Del(path string) (gorequest.Response, []byte, []error) {
	return c.requester.Delete(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).EndBytes()
}

func (c *GrClient) url() string {
	return c.config.RootUrl + "/" + c.config.ApiVersion
}

var Timeout int64

type Client struct {
	client *http.Client
	config *Config
}

func NewClient(cnf *Config) *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Second * time.Duration(Timeout),
		},
		config: cnf,
	}
}

func (c *Client) Get(path string, values url.Values, hdrs http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", c.url()+path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	for n, h := range hdrs {
		req.Header.Add(n, h[0])
	}

	req.URL.RawQuery = values.Encode()
	verboseHeader(req.Header, "GET-REQUEST")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	verboseHeader(resp.Header, "GET-RESPONSE")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) GetHeaders(path string, values url.Values, hdrs http.Header) (http.Header, error) {
	req, err := http.NewRequest("GET", c.url()+path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	for n, h := range hdrs {
		req.Header.Add(n, h[0])
	}

	req.URL.RawQuery = values.Encode()
	verboseHeader(req.Header, "HEADERS-REQUEST")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	verboseHeader(resp.Header, "HEADERS-RESPONSE")
	defer resp.Body.Close()
	ret := resp.Header
	return ret, nil
}

func (c *Client) GetAccessToken(pl []byte) (string, error) {
	hdrs := &http.Header{"Content-Type": {"application/json"}}
	body, err := c.post("/access_token", nil, hdrs, pl)
	if err != nil {
		return "", err
	}

	var m map[string]interface{}
	if err = json.Unmarshal(body, &m); err != nil {
		return "", err
	}

	token, found := m["access_token"]
	if !found {
		return "", fmt.Errorf("cannot find access token")
	}

	return fmt.Sprintf("%s", token), nil
}

func (c *Client) Del(path string, values url.Values) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.url()+path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	req.URL.RawQuery = values.Encode()
	verboseHeader(req.Header, "DEL-REQUEST")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	verboseHeader(resp.Header, "DEL-RESPONSE")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) url() string {
	return c.config.RootUrl + "/" + c.config.ApiVersion
}

func (c *Client) post(path string, v *url.Values, h *http.Header, pl []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.url()+path, bytes.NewBuffer(pl))
	if h != nil {
		for name, hdr := range *h {
			req.Header.Add(name, hdr[0])
		}
	}

	verboseRequest(req)
	verboseHeader(req.Header, "POST-REQUEST")

	if v != nil {
		values := *v
		req.URL.RawQuery = values.Encode()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	verboseHeader(resp.Header, "POST-RESPONSE")
	verboseResponse(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func verboseRequest(r *http.Request) {
	if d.Verbose {
		d.Info("[URL]", r.URL.String())
	}
}

func verboseResponse(r *http.Response) {
	if d.Verbose {
		d.Info("[STATUS]", r.Status)
	}
}

func verboseHeader(hdr http.Header, prefix string) {
	if d.Verbose {
		for n, h := range hdr {
			d.Info(fmt.Sprintf("[%s] %s: %s", prefix, n, h))
		}
	}
}
