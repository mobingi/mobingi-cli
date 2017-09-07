package sesha3

import "testing"

func TestNewSimpleSuccess(t *testing.T) {
	s := NewSimpleSuccess("test")
	if s.Status != "success" {
		t.Fatal("expected success")
	}

	if s.Description != "test" {
		t.Fatal("expected test")
	}
}

func TestNewSimpleError(t *testing.T) {
	e := NewSimpleError("test")
	if e.Status != "error" {
		t.Fatal("expected error")
	}

	if e.Description != "test" {
		t.Fatal("expected test")
	}

	if e.Trace == nil {
		t.Fatal("trace should not be nil")
	}
}
