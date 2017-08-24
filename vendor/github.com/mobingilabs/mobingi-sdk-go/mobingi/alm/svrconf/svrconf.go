package svrconf

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/pkg/errors"
)

type ServerConfigGetInput struct {
	StackId string
}

type ServerConfigUpdateEnvVarsInput struct {
	StackId string
	EnvVars string
}

type ServerConfigUpdateFilePathInput struct {
	StackId  string
	FilePath string
}

type cnf struct {
	session *session.Session
	client  client.HttpClient
}

func (s *cnf) Get(in *ServerConfigGetInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	ep := s.session.ApiEndpoint() + `/alm/serverconfig?stack_id=` + in.StackId
	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

func (s *cnf) UpdateEnvVars(in *ServerConfigUpdateEnvVarsInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	if in.EnvVars == "" {
		return nil, nil, errors.New("empty environment variable(s)")
	}

	env := s.buildEnvPayload(in.StackId, in.EnvVars)
	if env == "" {
		return nil, nil, errors.New("error in payload init")
	}

	return s.sendUpdates(in.StackId, env)
}

func (s *cnf) UpdateFilePath(in *ServerConfigUpdateFilePathInput) (*client.Response, []byte, error) {
	if in == nil {
		return nil, nil, errors.New("input cannot be nil")
	}

	if in.StackId == "" {
		return nil, nil, errors.New("stack id cannot be empty")
	}

	if in.FilePath == "" {
		return nil, nil, errors.New("empty filepath")
	}

	fp := s.buildFilePathPayload(in.StackId, in.FilePath)
	if fp == "" {
		return nil, nil, errors.New("error in payload init")
	}

	return s.sendUpdates(in.StackId, fp)
}

func (s *cnf) sendUpdates(id, in string) (*client.Response, []byte, error) {
	rm := json.RawMessage(in)
	payload, err := json.Marshal(&rm)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal failed")
	}

	ep := s.session.ApiEndpoint() + `/alm/serverconfig?stack_id=` + id
	req, err := http.NewRequest(http.MethodPut, ep, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.session.AccessToken)
	return s.client.Do(req)
}

func (s *cnf) buildEnvPayload(id, env string) string {
	cnt := 0
	payload := `{"stack_id":"` + id + `",`

	// check if delete all
	if env == "null" {
		payload += `"envvars":{}}`
		return payload
	}

	if env != "" {
		line := `"envvars":{`
		envs := strings.Split(env, ",")
		for i, s := range envs {
			kv := strings.Split(s, ":")
			if len(kv) == 2 {
				line += `"` + strings.TrimSpace(kv[0]) + `":"` + strings.TrimSpace(kv[1]) + `"`
				cnt += 1
			}

			if i < len(envs)-1 {
				line += `,`
			}
		}

		line += `}`
		payload += line
	}

	payload += `}`
	if cnt == 0 {
		return ""
	}

	return payload
}

func (s *cnf) buildFilePathPayload(id, fp string) string {
	return `{"stack_id":"` + id + `","filepath":"` + fp + `"}`
}

func New(s *session.Session) *cnf {
	if s == nil {
		return nil
	}

	var c client.HttpClient
	if s.Config.HttpClientConfig != nil {
		c = client.NewSimpleHttpClient(s.Config.HttpClientConfig)
	} else {
		c = client.NewSimpleHttpClient()
	}

	return &cnf{session: s, client: c}
}
