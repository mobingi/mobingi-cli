package credentials

import (
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mobingilabs/mocli/pkg/constants"
)

func SaveToken(token string) error {
	hd, _ := homedir.Dir()
	folder := filepath.Join(hd, constants.CRED_FOLDER)
	_ = os.Mkdir(folder, os.ModePerm)
	cred := filepath.Join(folder, constants.CRED_FILE)
	return ioutil.WriteFile(cred, []byte(token), 0644)
}

func GetToken() ([]byte, error) {
	hd, _ := homedir.Dir()
	folder := filepath.Join(hd, constants.CRED_FOLDER)
	cred := filepath.Join(folder, constants.CRED_FILE)
	return ioutil.ReadFile(cred)
}
