package debug

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mobingilabs/mocli/pkg/cli/confmap"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func isError(err interface{}) bool {
	var (
		derr  error
		valid bool
	)

	switch err.(type) {
	case string:
		if err != "" {
			derr = errors.WithStack(fmt.Errorf(fmt.Sprintf("%s", err)))
			Error(err)
			valid = true
		}
	case error:
		if err != nil {
			e, _ := err.(error)
			derr = errors.WithStack(e)
			Error(err)
			valid = true
		}
	}

	if valid {
		if viper.GetBool(confmap.ConfigKey("debug")) {
			// stack trace from 'errors'
			fmt.Printf("%+v\n", derr)
		}
	}

	return valid
}

func ErrorExit(err interface{}, code int) {
	if valid := isError(err); valid {
		os.Exit(code)
	}
}

func IsHttpSuccess(code int) bool {
	// only 2xx = OK
	if code >= 200 && code < 300 {
		return true
	}

	return false
}

func ResponseError(r gorequest.Response, b []byte) string {
	errcnt := 0
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err.Error()
	}

	var serr, cem string
	if r.StatusCode != 200 {
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
