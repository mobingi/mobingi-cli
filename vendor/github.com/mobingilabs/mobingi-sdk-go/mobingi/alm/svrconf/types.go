package svrconf

type ServerConfig struct {
	Code          string `json:"code"`
	CodeDir       string `json:"codeDir"`
	Image         string `json:"image"`
	Ports         []int  `json:"ports"`
	Updated       int64  `json:"updated"`
	Version       string `json:"version"`
	GitPrivateKey string `json:"gitPrivateKey"`
	GitReference  string `json:"gitReference"`
}
