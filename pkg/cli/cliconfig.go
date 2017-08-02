package cli

import (
	"encoding/json"

	"github.com/mobingilabs/mocli/pkg/iohelper"
	"github.com/mobingilabs/mocli/pkg/pretty"
)

type cliConfig struct {
	RunEnv  string `json:"runenv"`
	Verbose bool   `json:"verbose"`
}

func SetDefaultCliConfig(f string) error {
	defcfg := cliConfig{RunEnv: "prod"}
	contents, err := json.MarshalIndent(defcfg, "", pretty.Indent(3))
	if err != nil {
		return err
	}

	return iohelper.WriteToFile(f, contents)
}
