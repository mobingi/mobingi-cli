package debug

import (
	"testing"

	"github.com/pkg/errors"
)

func TestErrx(t *testing.T) {
	err := errors.New("test error")
	errx(false, err)
	errx(true, err)
}
