package constants

const (
	CRED_FOLDER   = ".mocli"      // folder name for config file(s), created in home folder
	CRED_FILE     = "credentials" // we store access token here
	REGTOKEN_FILE = "regtoken"    // we store Docker Registry token here

	PROD_API_BASE = "https://api.mobingi.com"            // production API base url
	DEV_API_BASE  = "https://apidev.mobingi.com"         // dev API base url
	PROD_REG_BASE = "https://registry.mobingi.com"       // production Docker Registry base url
	DEV_REG_BASE  = "https://dockereg2.labs.mobingi.com" // dev Docker Registry base url

	DOCKER_API_VER = "v2" // Docker API version
)
