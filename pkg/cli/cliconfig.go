package cli

import (
	"io/ioutil"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/pkg/iohelper"
	yaml "gopkg.in/yaml.v2"
)

type CliConfig struct {
	AccessToken string `yaml:"access_token"`
	RunEnv      string `yaml:"runenv"`
	Verbose     bool   `yaml:"verbose"`
}

func SetDefaultCliConfig() error {
	defcnf := CliConfig{RunEnv: "prod"}
	return defcnf.WriteToConfig()
}

func ReadCliConfig() *CliConfig {
	cnf := &CliConfig{}
	err := cnf.Reload()
	if err != nil {
		return nil
	}

	return cnf
}

func (c *CliConfig) WriteToConfig() error {
	contents, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return iohelper.WriteToFile(c.ConfigFile(), contents)
}

func (c *CliConfig) Reload() error {
	contents, err := ioutil.ReadFile(c.ConfigFile())
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(contents, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *CliConfig) ConfigFile() string {
	p := c.path()
	if p == "" {
		return p
	}

	cnf := filepath.Join(p, "."+BinName())
	return filepath.Join(cnf, "config")
}

func (c *CliConfig) path() string {
	var p string
	p, _ = homedir.Dir()
	return p
}
