package session

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	s1, _ := NewSession()
	if s1 == nil {
		t.Errorf("Expected non-nil session")
	}

	s2, _ := NewSession(&Config{})
	if s2 == nil {
		t.Errorf("Expected non-nil session")
	}

	s3, _ := NewSession(&Config{
		ClientId:     "clientid",
		ClientSecret: "clientsecret",
	})

	if s3.Config.ClientId != "clientid" {
		t.Errorf("Expected value 'clientid', got %s", s3.Config.ClientId)
	}

	s4, _ := NewSession(&Config{ApiVersion: 3})
	if s4.ApiEndpoint() != "https://alm.mobingi.com/v3" {
		t.Errorf("Invalid api url")
	}
}
