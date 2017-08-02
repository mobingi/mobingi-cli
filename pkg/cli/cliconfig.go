package cli

import (
	"encoding/json"

	"github.com/mobingilabs/mocli/pkg/iohelper"
	"github.com/mobingilabs/mocli/pkg/pretty"
)

type CliConfig struct {
	RunEnv  string `json:"runenv"`
	Verbose bool   `json:"verbose"`
}

func SetDefaultCliConfig(f string) error {
	defcfg := CliConfig{RunEnv: "prod"}
	contents, err := json.MarshalIndent(defcfg, "", pretty.Indent(3))
	if err != nil {
		return err
	}

	return iohelper.WriteToFile(f, contents)
}
