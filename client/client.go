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
)

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

func (c *Client) AuthGet(path string) ([]byte, error) {
	ah := c.authHdr()
	return c.get(path, nil, &ah)
}

func (c *Client) AuthPut(path string, pl []byte) ([]byte, error) {
	ah := c.authHdr()
	ah.Add("Content-Type", "application/json")
	return c.put(path, nil, &ah, pl)
}

func (c *Client) AuthDel(path string) ([]byte, error) {
	ah := c.authHdr()
	return c.del(path, nil, &ah)
}

func (c *Client) url() string {
	return c.config.RootUrl + "/" + c.config.ApiVersion
}

func (c *Client) hdr(path string, v *url.Values, h *http.Header) (http.Header, error) {
	req, err := http.NewRequest("GET", c.url()+path, nil)
	if err != nil {
		return nil, err
	}

	req = c.initReq(req, v, h)
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
	if err != nil {
		return nil, err
	}

	req = c.initReq(req, v, h)
	return c.send(req)
}

func (c *Client) post(path string, v *url.Values, h *http.Header, pl []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.url()+path, bytes.NewBuffer(pl))
	if err != nil {
		return nil, err
	}

	req = c.initReq(req, v, h)
	return c.send(req)
}

func (c *Client) put(path string, v *url.Values, h *http.Header, pl []byte) ([]byte, error) {
	req, err := http.NewRequest("PUT", c.url()+path, bytes.NewBuffer(pl))
	if err != nil {
		return nil, err
	}

	req = c.initReq(req, v, h)
	return c.send(req)
}

func (c *Client) del(path string, v *url.Values, h *http.Header) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.url()+path, nil)
	if err != nil {
		return nil, err
	}

	req = c.initReq(req, v, h)
	return c.send(req)
}

func (c *Client) send(r *http.Request) ([]byte, error) {
	verboseHeader(r.Header, fmt.Sprintf("%s-REQUEST", r.Method))

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	verboseHeader(resp.Header, fmt.Sprintf("%s-RESPONSE", r.Method))
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

func (c *Client) authHdr() http.Header {
	return http.Header{"Authorization": {"Bearer " + c.config.AccessToken}}
}

func (c *Client) initReq(r *http.Request, v *url.Values, h *http.Header) *http.Request {
	if h != nil {
		for name, hdr := range *h {
			r.Header.Add(name, hdr[0])
		}
	}

	if v != nil {
		values := *v
		r.URL.RawQuery = values.Encode()
	}

	verboseRequest(r)
	return r
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
