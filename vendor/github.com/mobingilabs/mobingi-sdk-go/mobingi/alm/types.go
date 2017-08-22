package alm

type Configuration struct {
	AWS                 string      `json:"AWS,omitempty"`
	AWSAccountName      string      `json:"AWS_ACCOUNT_NAME,omitempty"`
	AssociatePublicIp   string      `json:"AssociatePublicIP,omitempty"`
	ELBOpen443Port      string      `json:"ELBOpen443Port,omitempty"`
	ELBOpen80Port       string      `json:"ELBOpen80Port,omitempty"`
	SpotInstanceMaxSize int         `json:"SpotInstanceMaxSize,omitempty"`
	SpotInstanceMinSize int         `json:"SpotInstanceMinSize,omitempty"`
	SpotPrice           string      `json:"SpotPrice,omitempty"`
	Architecture        string      `json:"architecture,omitempty"`
	Code                string      `json:"code,omitempty"`
	Image               string      `json:"image,omitempty"`
	Max                 interface{} `json:"max,omitempty"`
	MaxOrigin           interface{} `json:"maxOrigin,omitempty"`
	Min                 interface{} `json:"min,omitempty"`
	MinOrigin           interface{} `json:"minOrigin,omitempty"`
	Nickname            string      `json:"nickname,omitempty"`
	Region              string      `json:"region,omitempty"`
	Type                string      `json:"type,omitempty"`
}

type StackOutput struct {
	// list
	Description string `json:"Description,omitempty"`
	OutputKey   string `json:"OutputKey,omitempty"`
	OutputValue string `json:"OutputValue,omitempty"`
	// describe
	Address                     string `json:"Address,omitempty"`
	DBAddress                   string `json:"DBAddress,omitempty"`
	DBPort                      string `json:"DBPort,omitempty"`
	DBSlave1Address             string `json:"DBSlave1Address,omitempty"`
	DBSlave2Address             string `json:"DBSlave2Address,omitempty"`
	DBSlave3Address             string `json:"DBSlave3Address,omitempty"`
	DBSlave4Address             string `json:"DBSlave4Address,omitempty"`
	DBSlave5Address             string `json:"DBSlave5Address,omitempty"`
	MemcachedEndPointAddress    string `json:"MemcachedEndPointAddress,omitempty"`
	MemcachedEndPointPort       string `json:"MemcachedEndPointPort,omitempty"`
	NATInstance                 string `json:"NATInstance,omitempty"`
	RedisPrimaryEndPointAddress string `json:"RedisPrimaryEndPointAddress,omitempty"`
	RedisPrimaryEndPointPort    string `json:"RedisPrimaryEndPointPort,omitempty"`
	RedisReadEndPointAddresses  string `json:"RedisReadEndPointAddresses,omitempty"`
	RedisReadEndPointPorts      string `json:"RedisReadEndPointPorts,omitempty"`
}

type ListStack struct {
	AuthToken     string        `json:"auth_token,omitempty"`
	Configuration Configuration `json:"configuration,omitempty"`
	CreateTime    string        `json:"create_time,omitempty"`
	Nickname      string        `json:"nickname,omitempty"`
	StackId       string        `json:"stack_id,omitempty"`
	StackOutputs  []StackOutput `json:"stack_outputs,omitempty"`
	StackStatus   string        `json:"stack_status,omitempty"`
	UserId        string        `json:"user_id,omitempty"`
}
