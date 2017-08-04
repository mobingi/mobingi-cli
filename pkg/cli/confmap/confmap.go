package confmap

var cfm = map[string]string{
	"token":   "access_token",
	"url":     "api_url",
	"rurl":    "registry_url",
	"apiver":  "api_version",
	"indent":  "indent",
	"timeout": "timeout",
	"verbose": "verbose",
	"debug":   "debug",
}

func ConfigKey(flag string) string {
	return cfm[flag]
}
