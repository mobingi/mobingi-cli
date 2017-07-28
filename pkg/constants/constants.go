package constants

const (
	CRED_FILE     = "credentials" // we store access token here
	REGTOKEN_FILE = "regtoken"    // we store Docker Registry token here

	PROD_API_BASE = "https://api.mobingi.com"      // production API base url
	PROD_REG_BASE = "https://registry.mobingi.com" // production Docker Registry base url

	QA_API_BASE = "https://apiqa.mobingi.com" // test API base url
	QA_REG_BASE = PROD_REG_BASE               // test Docker Registry base url

	DEV_API_BASE = "https://apidev.mobingi.com"         // dev API base url
	DEV_REG_BASE = "https://dockereg2.labs.mobingi.com" // dev Docker Registry base url

	DOCKER_API_VER = "v2" // Docker API version
)
