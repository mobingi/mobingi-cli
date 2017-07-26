package check

import (
	"encoding/json"
	"fmt"
	"os"

	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/parnurzeal/gorequest"
)

func isError(err interface{}) bool {
	valid := false
	switch err.(type) {
	case string:
		if err != "" {
			d.Error(err)
			valid = true
		}
	case error:
		if err != nil {
			d.Error(err)
			valid = true
		}
	case []error:
		s, ok := err.([]error)
		if ok {
			if len(s) > 0 {
				d.Error(s)
				valid = true
			}
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
