package svrconf

type ServerConfig struct {
	Code          string `json:"code,omitempty"`
	CodeDir       string `json:"codeDir,omitempty"`
	Image         string `json:"image,omitempty"`
	Ports         []int  `json:"ports,omitempty"`
	Updated       int64  `json:"updated,omitempty"`
	Version       string `json:"version,omitempty"`
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`
	GitReference  string `json:"gitReference,omitempty"`
}
