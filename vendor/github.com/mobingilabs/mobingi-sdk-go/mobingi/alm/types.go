package alm

import "encoding/json"

// Changes:
//
// 2017-07-18:
//   - Max, MaxOrigin, Min, MinOrigin - changed to int (we still need to support old string)
type Configuration struct {
	// v3
	Description string          `json:"description,omitempty"`
	Label       string          `json:"label,omitempty"`
	Version     string          `json:"version,omitempty"`
	Vendor      json.RawMessage `json:"vendor,omitempty"`
	// v2
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
	StackOutputs  interface{}   `json:"stack_outputs,omitempty"`
	StackStatus   string        `json:"stack_status,omitempty"`
	UserId        string        `json:"user_id,omitempty"`
}

type State struct {
	Code string `json:"Code,omitempty"`
	Name string `json:"Name,omitempty"`
}

type Instance struct {
	AmiLaunchIndex        string      `json:"AmiLaunchIndex,omitempty"`
	Architecture          string      `json:"Architecture,omitempty"`
	BlockDeviceMappings   interface{} `json:"BlockDeviceMappings,omitempty"`
	ClientToken           string      `json:"ClientToken,omitempty"`
	EbsOptimized          bool        `json:"EbsOptimized,omitempty"`
	Hypervisor            string      `json:"Hypervisor,omitempty"`
	ImageId               string      `json:"ImageId,omitempty"`
	InstanceId            string      `json:"InstanceId,omitempty"`
	InstanceType          string      `json:"InstanceType,omitempty"`
	InstanceLifecycle     string      `json:"InstanceLifecycle,omitempty"`
	SpotInstanceRequestId string      `json:"SpotInstanceRequestId,omitempty"`
	KeyName               string      `json:"KeyName,omitempty"`
	LaunchTime            string      `json:"LaunchTime,omitempty"`
	Monitoring            interface{} `json:"Monitoring,omitempty"`
	NetworkInterfaces     interface{} `json:"NetworkInterfaces,omitempty"`
	Placement             interface{} `json:"Placement,omitempty"`
	PrivateDnsName        string      `json:"PrivateDnsName,omitempty"`
	PrivateIpAddress      string      `json:"PrivateIpAddress,omitempty"`
	ProductCodes          []string    `json:"ProductCodes,omitempty"`
	PublicDnsName         string      `json:"PublicDnsName,omitempty"`
	PublicIpAddress       string      `json:"PublicIpAddress,omitempty"`
	Reservation           interface{} `json:"Reservation,omitempty"`
	RootDeviceName        string      `json:"RootDeviceName,omitempty"`
	RootDeviceType        string      `json:"RootDeviceType,omitempty"`
	SecurityGroups        interface{} `json:"SecurityGroups,omitempty"`
	SourceDestCheck       bool        `json:"SourceDestCheck,omitempty"`
	State                 State       `json:"State,omitempty"`
	StateTransitionReason string      `json:"StateTransitionReason,omitempty"`
	SubnetId              string      `json:"SubnetId,omitempty"`
	Tags                  interface{} `json:"Tags,omitempty"`
	VirtualizationType    string      `json:"VirtualizationType,omitempty"`
	VpcId                 string      `json:"VpcId,omitempty"`
	EnaSupport            string      `json:"enaSupport,omitempty"`
}

type DescribeStack struct {
	AuthToken     string        `json:"auth_token,omitempty"`
	Configuration Configuration `json:"configuration,omitempty"`
	CreateTime    string        `json:"create_time,omitempty"`
	Instances     []Instance    `json:"Instances,omitempty"`
	Nickname      string        `json:"nickname,omitempty"`
	StackId       string        `json:"stack_id,omitempty"`
	StackOutputs  interface{}   `json:"stack_outputs,omitempty"`
	StackStatus   string        `json:"stack_status,omitempty"`
	UserId        string        `json:"user_id,omitempty"`
}

type AlmTemplateVersion struct {
	VersionId    string `json:"version_id,omitempty"`
	Latest       bool   `json:"latest,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	Size         string `json:"size,omitempty"`
}
