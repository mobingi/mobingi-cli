package credentials

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type VendorCredentials struct {
	Id           string `json:"id,omitempty"`
	Account      string `json:"account,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}

type AWSCredentials struct {
	Name   string `json:"AWSAccountName,omitempty"`
	KeyId  string `json:"AWSSecretKey,omitempty"`
	Secret string `json:"AWSSecretKeyId,omitempty"`
}

type AddVendorCredentials struct {
	Credentials AWSCredentials `json:"credentials,omitempty"`
}

type UserPass struct {
	Username string
	Password string
}

// EnsureInput asks for input in stdin if username or password are empty.
// It will return two booleans when input was done from stdin ([0] for
// username, [1] for password.
func (up *UserPass) EnsureInput(allowEmpty bool) ([2]bool, error) {
	var in [2]bool
	if up.Username == "" {
		up.Username = Username()
		in[0] = true
	}

	if up.Password == "" {
		up.Password = Password()
		in[1] = true
	}

	if !allowEmpty {
		if up.Username == "" {
			return in, fmt.Errorf("username cannot be empty")
		}

		if up.Password == "" {
			return in, fmt.Errorf("password cannot be empty")
		}
	}

	return in, nil
}

type ClientIdSecret struct {
	Id     string
	Secret string
}

func (c *ClientIdSecret) EnsureInput(allowEmpty bool) error {
	if c.Id == "" {
		c.Id = ClientId()
	}

	if c.Secret == "" {
		c.Secret = ClientSecret()
	}

	if !allowEmpty {
		if c.Id == "" {
			return fmt.Errorf("client id cannot be empty")
		}

		if c.Secret == "" {
			return fmt.Errorf("client secret cannot be empty")
		}
	}

	return nil
}

func Username() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	user, _ := reader.ReadString('\n')
	return strings.TrimSpace(user)
}

func Password() string {
	fmt.Print("Password: ")
	pass, _ := terminal.ReadPassword(int(syscall.Stdin))
	return string(pass)
}

func ClientId() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Client ID: ")
	id, _ := reader.ReadString('\n')
	return strings.TrimSpace(id)
}

func ClientSecret() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Client secret: ")
	secret, _ := reader.ReadString('\n')
	return strings.TrimSpace(secret)
}
