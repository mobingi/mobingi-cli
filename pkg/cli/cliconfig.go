package cli

import (
	"github.com/mobingilabs/mocli/pkg/iohelper"
	yaml "gopkg.in/yaml.v2"
)

type CliConfig struct {
	/*
		RunEnv  string `json:"runenv"`
		Verbose bool   `json:"verbose"`
	*/
	RunEnv  string `yaml:"runenv"`
	Verbose bool   `yaml:"verbose"`
}

func SetDefaultCliConfig(f string) error {
	defcfg := CliConfig{RunEnv: "prod"}
	// contents, err := json.MarshalIndent(defcfg, "", pretty.Indent(3))
	contents, err := yaml.Marshal(&defcfg)
	if err != nil {
		return err
	}

	return iohelper.WriteToFile(f, contents)
}
