package pretty

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestIndent(t *testing.T) {
	if i := Indent(4); i != "    " {
		t.Errorf("Expected four(4) whitespaces, got %v", i)
	}
}

type t1 struct {
	S string
}

type t2 struct {
	M   map[string]string
	I   int
	T1  t1
	Pt1 *t1
	St1 []t1
}

func TestJSON(t *testing.T) {
	mck := t2{
		M: map[string]string{"one": "1", "two": "2"},
		I: 100,
		T1: t1{
			S: "struct",
		},
		Pt1: &t1{
			S: "struct pointer",
		},
		St1: make([]t1, 0),
	}

	// test marshal
	log.Println("with marshal:")
	mck.St1 = append(mck.St1, t1{S: "hello"})
	s := JSON(mck, 2)
	fmt.Println(s)
	if len(s) == 0 {
		t.Errorf("Not expecting an empty result")
	}

	switch s[0] {
	case '{', '[':
	default:
		t.Errorf("Invalid json")
	}

	// test direct string
	log.Println("string input:")
	b, err := json.Marshal(mck)
	if err != nil {
		t.Errorf("Marshal failed: %#v", err)
	}

	fmt.Println(JSON(string(b), 2))
}
