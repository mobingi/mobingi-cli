package util

import (
	"net/http"
	"testing"
)

func TestIndent(t *testing.T) {
	if i := Indent(4); i != "    " {
		t.Errorf("Expected four(4) whitespaces, got %v", i)
	}
}

func TestResponseError(t *testing.T) {
	r := &http.Response{
		StatusCode: 401,
	}

	b := []byte(`{}`)
	if re := ResponseError(r, b); re == "" {
		t.Errorf("Expected error mesage, got %v", re)
	}

	r.StatusCode = 200
	if re := ResponseError(r, b); re != "" {
		t.Errorf("Expected empty mesage, got %v", re)
	}

	b = []byte(`{"code":"100","error":"error":"message":"error message"}`)
	if re := ResponseError(r, b); re == "" {
		t.Errorf("Expected error mesage, got %v", re)
	}
}
