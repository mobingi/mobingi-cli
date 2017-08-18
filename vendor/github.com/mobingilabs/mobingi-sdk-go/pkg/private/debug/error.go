package debug

import "os"

func ErrorExit(err interface{}, code int) {
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

	if valid {
		os.Exit(code)
	}
}
