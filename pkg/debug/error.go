package debug

import (
	"os"
)

func isError(err interface{}) bool {
	var valid bool
	switch err.(type) {
	case string:
		if err != "" {
			Error(err)
			valid = true
		}
	case error:
		if err != nil {
			Error(err)
			valid = true
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
