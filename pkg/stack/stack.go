package stack

type Configuration struct {
	AWS                 string `json:"AWS,omitempty"`
	AWSAccountName      string `json:"AWS_ACCOUNT_NAME,omitempty"`
	AssociatePublicIp   string `json:"AssociatePublicIP,omitempty"`
	ELBOpen443Port      string `json:"ELBOpen443Port,omitempty"`
	ELBOpen80Port       string `json:"ELBOpen80Port,omitempty"`
	SpotInstanceMaxSize string `json:"SpotInstanceMaxSize,omitempty"`
	SpotInstanceMinSize string `json:"SpotInstanceMinSize,omitempty"`
	SpotPrice           string `json:"SpotPrice,omitempty"`
	Architecture        string `json:"architecture,omitempty"`
	Code                string `json:"code,omitempty"`
	Image               string `json:"image,omitempty"`
	Max                 string `json:"max,omitempty"`
	MaxOrigin           string `json:"maxOrigin,omitempty"`
	Min                 string `json:"min,omitempty"`
	MinOrigin           string `json:"minOrigin,omitempty"`
	Nickname            string `json:"nickname,omitempty"`
	Region              string `json:"region,omitempty"`
	Type                string `json:"type,omitempty"`
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

type Ebs struct {
	AttachTime          string `json:"AttachTime,omitempty"`
	DeleteOnTermination bool   `json:"DeleteOnTermination,omitempty"`
	Status              string `json:"Status,omitempty"`
	VolumeId            string `json:"VolumeId,omitempty"`
}

type BlockDeviceMappings struct {
	DeviceName string `json:"DeviceName,omitempty"`
	Ebs        Ebs    `json:"Ebs,omitempty"`
}

type Monitoring struct {
	State string `json:"State,omitempty"`
}

type Association struct {
	IpOwnerId     string `json:"IpOwnerId,omitempty"`
	PublicDnsName string `json:"PublicDnsName,omitempty"`
	PublicIp      string `json:"PublicIp,omitempty"`
}

type Attachment struct {
	AttachTime          string `json:"AttachTime,omitempty"`
	AttachmentId        string `json:"AttachmentId,omitempty"`
	DeleteOnTermination bool   `json:"DeleteOnTermination,omitempty"`
	DeviceIndex         string `json:"DeviceIndex,omitempty"`
	Status              string `json:"Status,omitempty"`
}

type Group struct {
	GroupId   string `json:"GroupId,omitempty"`
	GroupName string `json:"GroupName,omitempty"`
}

type PrivateIpAddress struct {
	Association      Association `json:"Association,omitempty"`
	Primary          bool        `json:"Primary,omitempty"`
	PrivateDnsName   string      `json:"PrivateDnsName,omitempty"`
	PrivateIpAddress string      `json:"PrivateIpAddress,omitempty"`
}

type NetworkInterface struct {
	Association        Association        `json:"Association,omitempty"`
	Attachment         Attachment         `json:"Attachment,omitempty"`
	Description        string             `json:"Description,omitempty"`
	Groups             []Group            `json:"Groups,omitempty"`
	MacAddress         string             `json:"MacAddress,omitempty"`
	NetworkInterfaceId string             `json:"NetworkInterfaceId,omitempty"`
	OwnerId            string             `json:"OwnerId,omitempty"`
	PrivateDnsName     string             `json:"PrivateDnsName,omitempty"`
	PrivateIpAddress   string             `json:"PrivateIpAddress,omitempty"`
	PrivateIpAddresses []PrivateIpAddress `json:"PrivateIpAddresses,omitempty"`
	SourceDestCheck    bool               `json:"SourceDestCheck,omitempty"`
	Status             string             `json:"Status,omitempty"`
	SubnetId           string             `json:"SubnetId,omitempty"`
	VpcId              string             `json:"VpcId,omitempty"`
}

type Placement struct {
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
	GroupName        string `json:"GroupName,omitempty"`
	Tenancy          string `json:"Tenancy,omitempty"`
}

type Reservation struct {
	Groups        []Group `json:"Groups,omitempty"`
	OwnerId       string  `json:"OwnerId,omitempty"`
	RequesterId   string  `json:"RequesterId,omitempty"`
	ReservationId string  `json:"ReservationId,omitempty"`
}

type State struct {
	Code string `json:"Code,omitempty"`
	Name string `json:"Name,omitempty"`
}

type Tag struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Instance struct {
	AmiLaunchIndex        string                `json:"AmiLaunchIndex,omitempty"`
	Architecture          string                `json:"Architecture,omitempty"`
	BlockDeviceMappings   []BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	ClientToken           string                `json:"ClientToken,omitempty"`
	EbsOptimized          bool                  `json:"EbsOptimized,omitempty"`
	Hypervisor            string                `json:"Hypervisor,omitempty"`
	ImageId               string                `json:"ImageId,omitempty"`
	InstanceId            string                `json:"InstanceId,omitempty"`
	InstanceType          string                `json:"InstanceType,omitempty"`
	KeyName               string                `json:"KeyName,omitempty"`
	LaunchTime            string                `json:"LaunchTime,omitempty"`
	Monitoring            Monitoring            `json:"Monitoring,omitempty"`
	NetworkInterfaces     []NetworkInterface    `json:"NetworkInterfaces,omitempty"`
	Placement             Placement             `json:"Placement,omitempty"`
	PrivateDnsName        string                `json:"PrivateDnsName,omitempty"`
	PrivateIpAddress      string                `json:"PrivateIpAddress,omitempty"`
	ProductCodes          []string              `json:"ProductCodes,omitempty"`
	PublicDnsName         string                `json:"PublicDnsName,omitempty"`
	PublicIpAddress       string                `json:"PublicIpAddress,omitempty"`
	Reservation           Reservation           `json:"Reservation,omitempty"`
	RootDeviceName        string                `json:"RootDeviceName,omitempty"`
	RootDeviceType        string                `json:"RootDeviceType,omitempty"`
	SecurityGroups        []Group               `json:"SecurityGroups,omitempty"`
	SourceDestCheck       bool                  `json:"SourceDestCheck,omitempty"`
	State                 State                 `json:"State,omitempty"`
	StateTransitionReason string                `json:"StateTransitionReason,omitempty"`
	SubnetId              string                `json:"SubnetId,omitempty"`
	Tags                  []Tag                 `json:"Tags,omitempty"`
	VirtualizationType    string                `json:"VirtualizationType,omitempty"`
	VpcId                 string                `json:"VpcId,omitempty"`
	EnaSupport            string                `json:"enaSupport,omitempty"`
}

// Separate struct due to `describe` output's stack_output not being an array.
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

type DescribeStack struct {
	AuthToken     string        `json:"auth_token,omitempty"`
	Configuration Configuration `json:"configuration,omitempty"`
	CreateTime    string        `json:"create_time,omitempty"`
	Instances     []Instance    `json:"Instances,omitempty"`
	Nickname      string        `json:"nickname,omitempty"`
	StackId       string        `json:"stack_id,omitempty"`
	StackOutputs  StackOutput   `json:"stack_outputs,omitempty"`
	StackStatus   string        `json:"stack_status,omitempty"`
	UserId        string        `json:"user_id,omitempty"`
}
