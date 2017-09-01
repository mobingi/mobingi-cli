package nativestore

import dcred "github.com/docker/docker-credential-helpers/credentials"

func Set(lbl, url, user, secret string) error {
	cr := &dcred.Credentials{
		ServerURL: url,
		Username:  user,
		Secret:    secret,
	}

	dcred.SetCredsLabel(lbl)
	return ns.Add(cr)
}

func Get(lbl, url string) (string, string, error) {
	dcred.SetCredsLabel(lbl)
	return ns.Get(url)
}
