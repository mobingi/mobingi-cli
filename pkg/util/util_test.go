package util

import "testing"

func TestIndent(t *testing.T) {
	i := Indent(4)
	if i != "    " {
		t.Errorf("Expected four(4) whitespaces, got %v", i)
	}
}
