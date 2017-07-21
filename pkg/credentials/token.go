package credentials

import (
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/pkg/constants"
)

func CredFolder(create bool) string {
	hd, _ := homedir.Dir()
	folder := filepath.Join(hd, constants.CRED_FOLDER)
	if create {
		_ = os.Mkdir(folder, os.ModePerm)
	}

	return folder
}

func SaveToken(token string) error {
	folder := CredFolder(true)
	cred := filepath.Join(folder, constants.CRED_FILE)
	return ioutil.WriteFile(cred, []byte(token), 0644)
}

func GetToken() ([]byte, error) {
	folder := CredFolder(false)
	cred := filepath.Join(folder, constants.CRED_FILE)
	return ioutil.ReadFile(cred)
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
