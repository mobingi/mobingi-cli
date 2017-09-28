package debug

import "os"

func errx(trace bool, err interface{}) bool {
	var valid bool
	switch err.(type) {
	case string:
		if err != "" {
			valid = true
			if trace {
				ErrorD(err)
			} else {
				Error(err)
			}
		}
	case error:
		if err != nil {
			valid = true
			if trace {
				ErrorD(err)
			} else {
				Error(err)
			}
		}
	}

	return valid
}

func ErrorTraceExit(err interface{}, code int) {
	valid := errx(true, err)
	if valid {
		os.Exit(code)
	}
}

func ErrorExit(err interface{}, code int) {
	valid := errx(false, err)
	if valid {
		os.Exit(code)
	}
}
