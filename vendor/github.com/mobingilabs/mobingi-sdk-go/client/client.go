package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/pkg/errors"
)

type Config struct {
	Timeout time.Duration
	Verbose bool
	Logger  *log.Logger // stdout when nil, verbose should be true
}

type httpClient interface {
	Do(*http.Request) (*http.Response, []byte, error)
}

type simpleHttpClient struct {
	client *http.Client
	cnf    *Config
}

func (c *simpleHttpClient) Do(r *http.Request) (*http.Response, []byte, error) {
	if c.cnf.Verbose {
		if c.cnf.Logger == nil {
			debug.Info("[URL]", r.URL.String())
			debug.Info("[METHOD]", r.Method)
			for n, h := range r.Header {
				debug.Info(fmt.Sprintf("[REQUEST] %s: %s", n, h))
			}
		} else {
			c.cnf.Logger.Println("[URL]", r.URL.String())
			c.cnf.Logger.Println("[METHOD]", r.Method)
			for n, h := range r.Header {
				c.cnf.Logger.Println(fmt.Sprintf("[REQUEST] %s: %s", n, h))
			}
		}
	}

	var lctx context.Context
	var lcancel context.CancelFunc
	req := r
	if c.cnf.Timeout > 0 {
		ctx := req.Context()
		lctx, lcancel = context.WithTimeout(ctx, c.cnf.Timeout)
		req = r.WithContext(lctx)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return resp, nil, errors.Wrap(err, "do failed")
	}

	defer resp.Body.Close()
	if lcancel != nil {
		defer lcancel()
	}

	if c.cnf.Verbose {
		if c.cnf.Logger == nil {
			for n, h := range resp.Header {
				debug.Info(fmt.Sprintf("[RESPONSE] %s: %s", n, h))
			}

			debug.Info("[STATUS]", resp.Status)
		} else {
			for n, h := range resp.Header {
				c.cnf.Logger.Println(fmt.Sprintf("[RESPONSE] %s: %s", n, h))
			}

			c.cnf.Logger.Println("[STATUS]", resp.Status)
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, body, errors.Wrap(err, "readall failed")
	}

	return resp, body, nil
}

func NewSimpleHttpClient(cnf *Config) *simpleHttpClient {
	if cnf == nil {
		return &simpleHttpClient{
			client: http.DefaultClient,
			cnf: &Config{
				Timeout: time.Second * 120,
			},
		}
	}

	return &simpleHttpClient{
		client: http.DefaultClient,
		cnf:    cnf,
	}
}
