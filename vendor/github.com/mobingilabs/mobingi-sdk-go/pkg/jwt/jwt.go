package jwt

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/pkg/errors"
)

type WrapperClaims struct {
	Data map[string]interface{}
	jwt.StandardClaims
}

type jwtctx struct {
	Pub    []byte
	Prv    []byte
	PemPub string
	PemPrv string
	init   bool
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

// NewCtx initializes our jwt context. For now, it is expected that the pem files
// (private and public) are already in os.TempDir() + "/jwt/rsa/".
func NewCtx() (*jwtctx, error) {
	// TODO: transfer this to authd service
	tmpdir := os.TempDir() + "/jwt/rsa/"
	pempub := tmpdir + "token.pem.pub"
	pemprv := tmpdir + "token.pem"

	pubcache, err := ioutil.ReadFile(pempub)
	if err != nil {
		debug.Error(err)
		return nil, errors.Wrap(err, "pub readfile failed")
	}

	prvcache, err := ioutil.ReadFile(pemprv)
	if err != nil {
		debug.Error(err)
		return nil, errors.Wrap(err, "prv readfile failed")
	}

	ctx := jwtctx{
		PemPub: pempub,
		PemPrv: pemprv,
		Pub:    pubcache,
		Prv:    prvcache,
	}

	return &ctx, nil
}
