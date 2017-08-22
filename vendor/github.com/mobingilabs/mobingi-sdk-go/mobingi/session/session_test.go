package session

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	s1, _ := New()
	if s1 == nil {
		t.Errorf("Expected non-nil session")
	}

	s2, _ := New(&Config{})
	if s2 == nil {
		t.Errorf("Expected non-nil session")
	}

	s3, _ := New(&Config{
		ClientId:     "clientid",
		ClientSecret: "clientsecret",
	})

	if s3.Config.ClientId != "clientid" {
		t.Errorf("Expected value 'clientid', got %s", s3.Config.ClientId)
	}

	s4, _ := New(&Config{ApiVersion: 3})
	if s4.ApiEndpoint() != "https://alm.mobingi.com/v3" {
		t.Errorf("Invalid api url")
	}
}
