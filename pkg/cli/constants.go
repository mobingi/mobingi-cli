package cli

const (
	ConfigFileName = "config.yml" // our yaml-based config file name

	ProductionBaseApiUrl       = "https://api.mobingi.com"            // production API base url
	ProductionBaseRegistryUrl  = "https://registry.mobingi.com"       // production Docker Registry base url
	TestBaseApiUrl             = "https://apiqa.mobingi.com"          // test API base url
	TestBaseRegistryUrl        = ProductionBaseRegistryUrl            // test Docker Registry base url
	DevelopmentBaseApiUrl      = "https://apidev.mobingi.com"         // dev API base url
	DevelopmentBaseRegistryUrl = "https://dockereg2.labs.mobingi.com" // dev Docker Registry base url

	ApiVersion       = "v2" // Mobingi API version
	DockerApiVersion = "v2" // Docker API version
)
