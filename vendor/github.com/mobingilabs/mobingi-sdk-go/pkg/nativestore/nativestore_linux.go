package nativestore

import "github.com/docker/docker-credential-helpers/secretservice"

var ns = secretservice.Secretservice{}

func nativeStore() *secretservice.Secretservice {
	return nil
}
