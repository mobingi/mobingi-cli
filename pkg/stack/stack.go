package stack

import (
	"fmt"
	"io"
	"reflect"

	"github.com/mobingilabs/mocli/pkg/pretty"
)

type Configuration struct {
	AWS                 string `json:"AWS,omitempty"`
	AWSAccountName      string `json:"AWS_ACCOUNT_NAME,omitempty"`
	AssociatePublicIp   string `json:"AssociatePublicIP,omitempty"`
	ELBOpen443Port      string `json:"ELBOpen443Port,omitempty"`
	ELBOpen80Port       string `json:"ELBOpen80Port,omitempty"`
	SpotInstanceMaxSize int    `json:"SpotInstanceMaxSize,omitempty"`
	SpotInstanceMinSize int    `json:"SpotInstanceMinSize,omitempty"`
	SpotPrice           string `json:"SpotPrice,omitempty"`
	Architecture        string `json:"architecture,omitempty"`
	Code                string `json:"code,omitempty"`
	Image               string `json:"image,omitempty"`
	Max                 int    `json:"max,omitempty"`
	MaxOrigin           int    `json:"maxOrigin,omitempty"`
	Min                 int    `json:"min,omitempty"`
	MinOrigin           int    `json:"minOrigin,omitempty"`
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

// Workaround for inconsistencies in API output:
// When stack creation is still in progress, StackOutputs is a slice. Upon completion,
// it will be a struct. It will cause errors in Unmarshal.
type DescribeStack1 struct {
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

type DescribeStack2 struct {
	AuthToken     string        `json:"auth_token,omitempty"`
	Configuration Configuration `json:"configuration,omitempty"`
	CreateTime    string        `json:"create_time,omitempty"`
	Instances     []Instance    `json:"Instances,omitempty"`
	Nickname      string        `json:"nickname,omitempty"`
	StackId       string        `json:"stack_id,omitempty"`
	StackOutputs  []StackOutput `json:"stack_outputs,omitempty"`
	StackStatus   string        `json:"stack_status,omitempty"`
	UserId        string        `json:"user_id,omitempty"`
}

// PrintR prints the `field: value` of the input struct recursively. Recursion level `lvl` and `indent`
// are provided for indention in printing. For slices, we have to do an explicit type assertion
// to get the underlying slice from reflect.
func PrintR(w io.Writer, s interface{}, lvl, indent int) {
	pad := pretty.Indent(lvl * indent)
	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i).Name
		value := rv.Field(i).Interface()

		switch rv.Field(i).Kind() {
		case reflect.String:
			fmt.Fprintf(w, "%s%s: %s\n", pad, field, value)
		case reflect.Int32:
			fmt.Fprintf(w, "%s%s: %i\n", pad, field, value)
		case reflect.Struct:
			fmt.Fprintf(w, "%s[%s]\n", pad, field)
			v := rv.Field(i).Addr()
			PrintR(w, v.Interface(), lvl+1, indent)
		case reflect.Slice:
			fmt.Fprintf(w, "%s[%s]\n", pad, field)
			instances, ok := value.([]Instance)
			if ok && len(instances) > 0 {
				for _, slice := range instances {
					PrintR(w, &slice, lvl+1, indent)
					if len(instances) > 1 {
						fmt.Fprintf(w, "\n")
					}
				}

				break
			}

			mappings, ok := value.([]BlockDeviceMappings)
			if ok && len(mappings) > 0 {
				for _, slice := range mappings {
					PrintR(w, &slice, lvl+1, indent)
					if len(mappings) > 1 {
						fmt.Fprintf(w, "\n")
					}
				}

				break
			}

			networks, ok := value.([]NetworkInterface)
			if ok && len(networks) > 0 {
				for _, slice := range networks {
					PrintR(w, &slice, lvl+1, indent)
					if len(networks) > 1 {
						fmt.Fprintf(w, "\n")
					}
				}

				break
			}

			groups, ok := value.([]Group)
			if ok && len(groups) > 0 {
				for _, slice := range groups {
					PrintR(w, &slice, lvl+1, indent)
					if len(groups) > 1 {
						fmt.Fprintf(w, "\n")
					}
				}

				break
			}

			ipaddrs, ok := value.([]PrivateIpAddress)
			if ok && len(ipaddrs) > 0 {
				for _, slice := range ipaddrs {
					PrintR(w, &slice, lvl+1, indent)
					if len(ipaddrs) > 1 {
						fmt.Fprintf(w, "\n")
					}
				}

				break
			}

			tags, ok := value.([]Tag)
			if ok && len(tags) > 0 {
				for _, slice := range tags {
					PrintR(w, &slice, lvl+1, indent)
					if len(tags) > 1 {
						fmt.Fprintf(w, "\n")
					}
				}

				break
			}

			stackouts, ok := value.([]StackOutput)
			if ok && len(stackouts) > 0 {
				for _, slice := range stackouts {
					PrintR(w, &slice, lvl+1, indent)
					if len(stackouts) > 1 {
						fmt.Fprintf(w, "\n")
					}
				}

				break
			}

			// when slice type is not explicitly specified in our conversion
			fmt.Fprintf(w, "%s*** Not available ***\n", pad)
		}
	}
}
