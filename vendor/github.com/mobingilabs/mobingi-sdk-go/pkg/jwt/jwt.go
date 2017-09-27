package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/private"
	"github.com/pkg/errors"
)

var rsainit bool

func init() {
	tmpdir := os.TempDir() + "/sesha3/rsa/"
	debug.Info("tmp:", tmpdir)
	pub := tmpdir + "token.pem.pub"
	prv := tmpdir + "token.pem"

	// create dir if necessary
	if !private.Exists(tmpdir) {
		err := os.MkdirAll(tmpdir, 0700)
		if err != nil {
			return
		}
	}

	// create public and private pem files
	if !private.Exists(pub) || !private.Exists(prv) {
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return
		}

		privder := x509.MarshalPKCS1PrivateKey(priv)
		pubkey := priv.Public()
		pubder, err := x509.MarshalPKIXPublicKey(pubkey)
		if err != nil {
			return
		}

		pubblock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubder}
		pemblock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privder}
		pubfile, err := os.OpenFile(pub, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return
		}

		defer pubfile.Close()
		prvfile, err := os.OpenFile(prv, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return
		}

		defer prvfile.Close()
		pem.Encode(pubfile, pubblock)
		pem.Encode(prvfile, pemblock)
	}

	rsainit = true
}

type WrapperClaims struct {
	Data map[string]interface{}
	jwt.StandardClaims
}

type jwtctx struct {
	Pub    []byte
	Prv    []byte
	PemPub string
	PemPrv string
}

func (j *jwtctx) GenerateToken(data map[string]interface{}) (*jwt.Token, string, error) {
	var stoken string
	var claims WrapperClaims

	claims.Data = data
	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), claims)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.Prv)
	if err != nil {
		return token, stoken, errors.Wrap(err, "parse priv key from pem failed")
	}

	stoken, err = token.SignedString(key)
	if err != nil {
		return token, stoken, errors.Wrap(err, "signed string failed")
	}

	return token, stoken, nil
}

func (j *jwtctx) ParseToken(token string) (*jwt.Token, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.Pub)
	if err != nil {
		return nil, errors.Wrap(err, "ParseRSAPublicKeyFromPEM failed")
	}

	var claims WrapperClaims
	return jwt.ParseWithClaims(token, &claims, func(tk *jwt.Token) (interface{}, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tk.Header["alg"])
		}

		return key, nil
	})
}

func NewCtx() (*jwtctx, error) {
	if !rsainit {
		return nil, errors.New("failed in rsa init")
	}

	var ctx jwtctx
	var err error

	tmpdir := os.TempDir() + "/sesha3/rsa/"
	ctx.PemPub = tmpdir + "token.pem.pub"
	ctx.PemPrv = tmpdir + "token.pem"
	ctx.Pub, err = ioutil.ReadFile(ctx.PemPub)
	if err != nil {
		return nil, errors.Wrap(err, "read public pem failed")
	}

	ctx.Prv, err = ioutil.ReadFile(ctx.PemPrv)
	if err != nil {
		return nil, errors.Wrap(err, "read private pem failed")
	}

	return &ctx, nil
}
