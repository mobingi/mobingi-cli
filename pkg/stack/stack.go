package stack

type Configuration struct {
	Aws                 string `json:"AWS,omitempty"`
	AssociatePublicIp   string `json:"AssociatePublicIP,omitempty"`
	ELBOpen443Port      string `json:"ELBOpen443Port,omitempty"`
	ELBOpen80Port       string `json:"ELBOpen80Port,omitempty"`
	SpotInstanceMaxSize string `json:"SpotInstanceMaxSize,omitempty"`
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
	Description string `json:"Description,omitempty"`
	OutputKey   string `json:"OutputKey,omitempty"`
	OutputValue string `json:"OutputValue,omitempty"`
}

type Stack struct {
	AuthToken     string        `json:"auth_token,omitempty"`
	Configuration Configuration `json:"configuration,omitempty"`
	CreateTime    string        `json:"create_time,omitempty"`
	Nickname      string        `json:"nickname,omitempty"`
	StackId       string        `json:"stack_id,omitempty"`
	StackOutputs  []StackOutput `json:"stack_outputs,omitempty"`
	StackStatus   string        `json:"stack_status,omitempty"`
	UserId        string        `json:"user_id,omitempty"`
}
