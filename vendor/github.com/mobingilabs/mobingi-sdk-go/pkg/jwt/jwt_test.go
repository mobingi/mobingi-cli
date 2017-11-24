package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/private"
)

func initpem() {
	tmpdir := os.TempDir() + "/jwt/rsa/"
	pempub := tmpdir + "token.pem.pub"
	pemprv := tmpdir + "token.pem"

	// create dir if necessary
	if !private.Exists(tmpdir) {
		err := os.MkdirAll(tmpdir, 0700)
		if err != nil {
			debug.Error(err)
			return
		}
	}

	// create public and private pem files
	if !private.Exists(pempub) || !private.Exists(pemprv) {
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			debug.Error(err)
			return
		}

		privder := x509.MarshalPKCS1PrivateKey(priv)
		pubkey := priv.Public()
		pubder, err := x509.MarshalPKIXPublicKey(pubkey)
		if err != nil {
			debug.Error(err)
			return
		}

		pubblock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubder}
		pemblock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privder}
		pubfile, err := os.OpenFile(pempub, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			debug.Error(err)
			return
		}

		defer pubfile.Close()
		prvfile, err := os.OpenFile(pemprv, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			debug.Error(err)
			return
		}

		defer prvfile.Close()
		err = pem.Encode(pubfile, pubblock)
		if err != nil {
			debug.Error(err)
			return
		}

		err = pem.Encode(prvfile, pemblock)
		if err != nil {
			debug.Error(err)
			return
		}
	}
}

func TestNewCtx(t *testing.T) {
	initpem()
	ctx, err := NewCtx()
	if err != nil {
		t.Fatal(err)
	}

	if ctx == nil {
		t.Fatal("should not be nil")
	}
}

func TestGenerateToken(t *testing.T) {
	initpem()
	ctx, _ := NewCtx()
	claims := make(map[string]interface{})
	claims["username"] = "user"
	token, stoken, _ := ctx.GenerateToken(nil)
	if token == nil {
		t.Fatal("should not be nil")
	}

	log.Println(token, stoken)
}

func TestParseToken(t *testing.T) {
	initpem()
	ctx, _ := NewCtx()
	claims := make(map[string]interface{})
	claims["username"] = "user"
	token, stoken, _ := ctx.GenerateToken(claims)
	if token == nil {
		t.Fatal("should not be nil")
	}

	pt, err := ctx.ParseToken(stoken)
	if err != nil {
		t.Fatal("should succeed; got:", err)
	}

	if !pt.Valid {
		t.Fatal("should be a valid token")
	}

	nc := pt.Claims.(*WrapperClaims)
	u, ok := nc.Data["username"]
	if !ok {
		t.Fatal("should have a username entry")
	}

	if fmt.Sprintf("%s", u) != "user" {
		t.Fatal("should be user")
	}

	log.Println(stoken)
}
