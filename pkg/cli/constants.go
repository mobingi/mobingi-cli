package cli

const (
	ConfigFileName = "config.yml" // our yaml-based config file name

	// v2 base
	ProductionBaseApiUrl       = "https://api.mobingi.com"            // production API base url
	ProductionBaseRegistryUrl  = "https://registry.mobingi.com"       // production Docker Registry base url
	TestBaseApiUrl             = "https://apiqa.mobingi.com"          // test API base url
	TestBaseRegistryUrl        = ProductionBaseRegistryUrl            // test Docker Registry base url
	DevelopmentBaseApiUrl      = "https://apidev.mobingi.com"         // dev API base url
	DevelopmentBaseRegistryUrl = "https://dockereg2.labs.mobingi.com" // dev Docker Registry base url

	// v3 base
	AlmBaseApiUrl      = "https://alm.mobingi.com"     // production API base url for ALM
	WaveBaseApiUrl     = "https://wave.mobingi.com"    // production API base url for Wave
	AlmDevBaseApiUrl   = "https://almdev.mobingi.com"  // dev API base url for ALM
	WaveDevBaseApiUrl  = "https://wavedev.mobingi.com" // dev API base url for Wave
	AlmTestBaseApiUrl  = "https://almqa.mobingi.com"   // test API base url for ALM
	WaveTestBaseApiUrl = "https://waveqamobingi.com"   // test API base url for Wave

	ApiVersion       = "v2" // Mobingi API version
	DockerApiVersion = "v2" // Docker API version
)
