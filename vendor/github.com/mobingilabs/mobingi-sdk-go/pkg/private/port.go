package private

import (
	"net"

	"github.com/pkg/errors"
)

func GetFreePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return -1, errors.Wrap(err, "listen failed")
	}

	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
