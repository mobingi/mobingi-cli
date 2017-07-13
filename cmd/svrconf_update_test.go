package cmd

import (
	"encoding/json"
	"testing"
)

func TestBuildPayload(t *testing.T) {
	v := buildPayload("sid", "null")
	m := json.RawMessage(v)
	_, err := json.Marshal(&m)
	if err != nil {
		t.Errorf("Expected proper marshaling, got error %v", err)
	}
}
