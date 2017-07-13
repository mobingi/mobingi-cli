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

	v = buildPayload("sid", "KEY:value")
	m = json.RawMessage(v)
	_, err = json.Marshal(&m)
	if err != nil {
		t.Errorf("Expected proper marshaling, got error %v", err)
	}

	v = buildPayload("sid", "KEY1:value1,KEY2:value2")
	m = json.RawMessage(v)
	_, err = json.Marshal(&m)
	if err != nil {
		t.Errorf("Expected proper marshaling, got error %v", err)
	}

	v = buildPayload("sid", "")
	if v != "" {
		t.Errorf("Expected an empty string, got %s", v)
	}

	v = buildPayload("sid", "KEY1,value1,KEY2,value2")
	if v != "" {
		t.Errorf("Expected an empty string, got %s", v)
	}
}
