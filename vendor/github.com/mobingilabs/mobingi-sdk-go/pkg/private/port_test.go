package private

import "testing"

func TestGetFreePort(t *testing.T) {
	p, _ := GetFreePort()
	if p <= 0 {
		t.Fatal("expecting > 0 port, got:", p)
	}
}
