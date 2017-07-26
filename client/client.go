package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mobingilabs/mocli/pkg/check"
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

func (c *Client) GetStacks() ([]byte, error) {
	hdr := &http.Header{"Authorization": {"Bearer " + c.config.AccessToken}}
	return c.get("/alm/stack", nil, hdr)
}

func (c *Client) GetStack(id string) ([]byte, error) {
	hdr := &http.Header{"Authorization": {"Bearer " + c.config.AccessToken}}
	return c.get("/alm/stack/"+id, nil, hdr)
}

func (c *Client) GetTagDigest(path string) (string, error) {
	var (
		digest string
	)

	hdrs := &http.Header{
		"Authorization": {"Bearer " + c.config.AccessToken},
		"Accept":        {"application/vnd.docker.distribution.manifest.v2+json"},
	}

	h, err := c.hdr(path, nil, hdrs)
	if err != nil {
		return digest, err
	}

	for name, hdr := range h {
		if name == "Etag" {
			digest = hdr[0]
			digest = strings.TrimSuffix(strings.TrimPrefix(digest, "\""), "\"")
		}
	}

	if digest == "" {
		return digest, fmt.Errorf("digest not found")
	}

	return digest, nil
}

func (c *Client) GetAccessToken(pl []byte) (string, error) {
	var (
		token string
		m     map[string]interface{}
	)

	hdrs := &http.Header{"Content-Type": {"application/json"}}
	body, err := c.post("/access_token", nil, hdrs, pl)
	if err != nil {
		return token, err
	}

	if err = json.Unmarshal(body, &m); err != nil {
		return token, err
	}

	t, found := m["access_token"]
	if !found {
		return token, fmt.Errorf("cannot find access token")
	}

	token = fmt.Sprintf("%s", t)
	return token, nil
}

func (c *Client) GetRegistryCatalog() ([]byte, error) {
	hdrs := &http.Header{"Authorization": {"Bearer " + c.config.AccessToken}}
	return c.get("/_catalog", nil, hdrs)
}

func (c *Client) GetRegistryTags(path string) ([]byte, error) {
	hdrs := &http.Header{"Authorization": {"Bearer " + c.config.AccessToken}}
	return c.get(path, nil, hdrs)
}

func (c *Client) GetRegistryTagManifest(path string) ([]byte, error) {
	hdrs := &http.Header{"Authorization": {"Bearer " + c.config.AccessToken}}
	return c.get(path, nil, hdrs)
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

func (c *Client) hdr(path string, v *url.Values, h *http.Header) (http.Header, error) {
	req, err := http.NewRequest("GET", c.url()+path, nil)
	if err != nil {
		return nil, err
	}

	if h != nil {
		for name, hdr := range *h {
			req.Header.Add(name, hdr[0])
		}
	}

	if v != nil {
		values := *v
		req.URL.RawQuery = values.Encode()
	}

	verboseRequest(req)
	verboseHeader(req.Header, "HEADERS-REQUEST")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	verboseHeader(resp.Header, "HEADERS-RESPONSE")
	verboseResponse(resp)
	ret := resp.Header
	return ret, nil
}

func (c *Client) get(path string, v *url.Values, h *http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", c.url()+path, nil)
	if h != nil {
		for name, hdr := range *h {
			req.Header.Add(name, hdr[0])
		}
	}

	if v != nil {
		values := *v
		req.URL.RawQuery = values.Encode()
	}

	verboseRequest(req)
	verboseHeader(req.Header, "GET-REQUEST")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	verboseHeader(resp.Header, "GET-RESPONSE")
	verboseResponse(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := respError(resp, body)
	if re != "" {
		return body, fmt.Errorf(re)
	}

	return body, nil
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

func respError(r *http.Response, b []byte) string {
	var (
		errcnt int
		m      map[string]interface{}
		serr   string
		cem    string
	)

	err := json.Unmarshal(b, &m)
	if err != nil {
		// considered success; our expected error format
		// is marshallable to 'm'
		return serr
	}

	if !check.IsHttpSuccess(r.StatusCode) {
		serr = serr + "[" + r.Status + "]"
	}

	// these three should be present to be considered error
	if c, found := m["code"]; found {
		cem = cem + "[" + fmt.Sprintf("%s", c) + "]"
		errcnt += 1
	}

	if e, found := m["error"]; found {
		cem = cem + fmt.Sprintf(" %s:", e)
		errcnt += 1
	}

	if s, found := m["message"]; found {
		cem = cem + fmt.Sprintf(" %s", s)
		errcnt += 1
	}

	if errcnt == 3 {
		serr = serr + cem
	}

	return serr
}
