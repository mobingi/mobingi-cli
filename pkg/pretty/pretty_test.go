package pretty

import (
	"testing"
)

func TestIndent(t *testing.T) {
	if i := Indent(4); i != "    " {
		t.Errorf("Expected four(4) whitespaces, got %v", i)
	}
}
