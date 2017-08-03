package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mobingilabs/mocli/client/timeout"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
)

// var Timeout int64 = 120 // same default value with cmdline flag (seconds)

type setreq struct {
	values *url.Values           // when not nil, we populate raw query
	header *http.Header          // when not nil, we add to headers
	basic  *credentials.UserPass // when not nil, we set basic auth
}

type Client struct {
	client *http.Client // our http client
	config *Config      // client configuration(s)
}

func NewClient(cnf *Config) *Client {
	return &Client{
		client: &http.Client{},
		config: cnf,
	}
}

func (c *Client) GetTagDigest(path string) (string, error) {
	var digest string
	hdrs := &http.Header{
		"Authorization": {"Bearer " + c.config.AccessToken},
		"Accept":        {"application/vnd.docker.distribution.manifest.v2+json"},
	}

	h, err := c.hdr(path, &setreq{header: hdrs})
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
	body, err := c.post("/access_token", &setreq{header: hdrs}, pl)
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

func (c *Client) BasicAuthGet(path, user, pass string, v *url.Values) ([]byte, error) {
	return c.get(
		path,
		&setreq{
			values: v,
			basic: &credentials.UserPass{
				Username: user,
				Password: pass,
			},
		},
	)
}

func (c *Client) AuthGet(path string) ([]byte, error) {
	ah := c.authHdr()
	return c.get(path, &setreq{header: &ah})
}

func (c *Client) AuthPost(path string, pl []byte) ([]byte, error) {
	ah := c.authHdr()
	ah.Add("Content-Type", "application/json")
	return c.post(path, &setreq{header: &ah}, pl)
}

func (c *Client) AuthPut(path string, pl []byte) ([]byte, error) {
	ah := c.authHdr()
	ah.Add("Content-Type", "application/json")
	return c.put(path, &setreq{header: &ah}, pl)
}

func (c *Client) AuthDel(path string) ([]byte, error) {
	ah := c.authHdr()
	return c.del(path, &setreq{header: &ah})
}

func (c *Client) url() string {
	return c.config.RootUrl + "/" + c.config.ApiVersion
}

func (c *Client) hdr(path string, p *setreq) (http.Header, error) {
	req, err := http.NewRequest(http.MethodGet, c.url()+path, nil)
	if err != nil {
		return nil, err
	}

	var cancel context.CancelFunc
	req, cancel = c.initReq(req, p)
	defer cancel()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	verboseResponse(resp)
	ret := resp.Header
	return ret, nil
}

func (c *Client) get(path string, p *setreq) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.url()+path, nil)
	if err != nil {
		return nil, err
	}

	var cancel context.CancelFunc
	req, cancel = c.initReq(req, p)
	return c.send(req, cancel)
}

func (c *Client) post(path string, p *setreq, pl []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.url()+path, bytes.NewBuffer(pl))
	if err != nil {
		return nil, err
	}

	var cancel context.CancelFunc
	req, cancel = c.initReq(req, p)
	return c.send(req, cancel)
}

func (c *Client) put(path string, p *setreq, pl []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPut, c.url()+path, bytes.NewBuffer(pl))
	if err != nil {
		return nil, err
	}

	var cancel context.CancelFunc
	req, cancel = c.initReq(req, p)
	return c.send(req, cancel)
}

func (c *Client) del(path string, p *setreq) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, c.url()+path, nil)
	if err != nil {
		return nil, err
	}

	var cancel context.CancelFunc
	req, cancel = c.initReq(req, p)
	return c.send(req, cancel)
}

func (c *Client) send(r *http.Request, cancel context.CancelFunc) ([]byte, error) {
	if cancel != nil {
		defer cancel()
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
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

func (c *Client) initReq(r *http.Request, p *setreq) (*http.Request, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(timeout.Timeout))
	r = r.WithContext(ctx)

	if p.header != nil {
		for name, hdr := range *p.header {
			r.Header.Add(name, hdr[0])
		}
	}

	if p.values != nil {
		values := *p.values
		r.URL.RawQuery = values.Encode()
	}

	if p.basic != nil {
		r.SetBasicAuth(p.basic.Username, p.basic.Password)
	}

	c.verboseRequest(r)
	return r, cancel
}

func (c *Client) verboseRequest(r *http.Request) {
	if d.Verbose {
		d.Info("[URL]", r.URL.String())
		d.Info("[METHOD]", r.Method)
		for n, h := range r.Header {
			d.Info(fmt.Sprintf("[REQUEST] %s: %s", n, h))
		}
	}
}

func verboseResponse(r *http.Response) {
	if d.Verbose {
		for n, h := range r.Header {
			d.Info(fmt.Sprintf("[RESPONSE] %s: %s", n, h))
		}

		d.Info("[STATUS]", r.Status)
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
