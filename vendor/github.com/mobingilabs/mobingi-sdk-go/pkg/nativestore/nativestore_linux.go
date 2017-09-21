package nativestore

import "github.com/docker/docker-credential-helpers/pass"

//passphrase of gpg2 is empty
//pass setting (below commands)
//$ gpg2 --gen-key 
//$ pass init [gpgID]
var ns = pass.Pass{}
