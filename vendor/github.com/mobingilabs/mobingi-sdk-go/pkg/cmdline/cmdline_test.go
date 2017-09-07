package cmdline

import "testing"

func TestArgs0(t *testing.T) {
	if Args0() == "" {
		t.Fatal("expected a name")
	}
}

func TestDir(t *testing.T) {
	if Dir() == "" {
		t.Fatal("expected a dir")
	}
}
