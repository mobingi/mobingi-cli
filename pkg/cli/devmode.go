package cli

var devMode bool

func DevMode() *bool {
	return &devMode
}

func IsDevMode() bool {
	return devMode
}
