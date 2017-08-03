package credentials

import (
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/constants"
)

func CredFolder(create bool) string {
	hd, _ := homedir.Dir()
	folder := filepath.Join(hd, "."+cli.BinName())
	if create {
		_ = os.Mkdir(folder, os.ModePerm)
	}

	return folder
}

func SaveRegistryToken(token string) error {
	folder := CredFolder(true)
	rf := filepath.Join(folder, constants.REGTOKEN_FILE)
	return ioutil.WriteFile(rf, []byte(token), 0644)
}

func GetRegistryToken() ([]byte, error) {
	folder := CredFolder(false)
	rf := filepath.Join(folder, constants.REGTOKEN_FILE)
	return ioutil.ReadFile(rf)
}
