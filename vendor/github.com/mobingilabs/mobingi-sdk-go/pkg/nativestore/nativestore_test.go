package nativestore

import (
	"log"
	"testing"
)

func TestSetGet(t *testing.T) {
	url := "https://github.com/mobingilabs/mobingi-sdk-go"
	Set("mobingi-sdk-go", url, "user", "password")
	user, secret, err := Get("mobingi-sdk-go", url)
	if err == nil {
		if user != "user" {
			t.Errorf("Expecting user, got %s", user)
		}

		if secret != "password" {
			t.Errorf("Expecting password, got %s", secret)
		}
	} else {
		log.Println("got error:", err)
	}
}
