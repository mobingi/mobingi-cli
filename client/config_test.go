package client

import (
	"testing"
)

func TestNewApiConfig(t *testing.T) {
	c := NewApiConfig(nil)
	if c != nil {
		t.Errorf("Expected a nil config")
	}
}
