package jwt

import (
	"fmt"
	"log"
	"testing"
)

func TestNewCtx(t *testing.T) {
	ctx, err := NewCtx()
	if err != nil {
		t.Fatal(err)
	}

	if ctx == nil {
		t.Fatal("should not be nil")
	}
}

func TestGenerateToken(t *testing.T) {
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
}
