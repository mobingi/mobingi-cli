package cli

import (
	"io/ioutil"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/client/timeout"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/iohelper"
	"github.com/mobingilabs/mocli/pkg/pretty"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// CliConfig is the object representation of our config file. The field tags for YAML marshaling and
// unmarshaling match the defined cli constants with prefix 'Config'.
type CliConfig struct {
	AccessToken string `yaml:"access_token"`
	// RunEnv          string `yaml:"run_env"`
	BaseApiUrl      string `yaml:"api_url"`
	BaseRegistryUrl string `yaml:"registry_url"`
	ApiVersion      string `yaml:"api_version"`
	Indent          int    `yaml:"indent"`
	Timeout         int64  `yaml:"timeout"`
	Verbose         bool   `yaml:"verbose"`
	Debug           bool   `yaml:"debug"`
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
	return filepath.Join(cnf, ConfigFileName)
}

func (c *CliConfig) path() string {
	var p string
	p, _ = homedir.Dir()
	return p
}

func SetDefaultCliConfig() error {
	defcnf := CliConfig{
		BaseApiUrl:      ProductionBaseApiUrl,
		BaseRegistryUrl: ProductionBaseRegistryUrl,
		ApiVersion:      ApiVersion,
		Indent:          pretty.Pad,
		Timeout:         timeout.Timeout,
	}

	err := defcnf.WriteToConfig()
	if err != nil {
		d.Error(err)
	}
	// check.ErrorExit(err, 1)

	return viper.ReadInConfig()
}

func ReadCliConfig() *CliConfig {
	cnf := &CliConfig{}
	err := cnf.Reload()
	if err != nil {
		return nil
	}

	return cnf
}
