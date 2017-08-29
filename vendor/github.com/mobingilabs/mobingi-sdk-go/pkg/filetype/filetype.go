package filetype

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

func IsJSON(in string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(in), &js) == nil
}

func IsYAML(in string) bool {
	var y yaml.MapSlice
	return yaml.Unmarshal([]byte(in), &y) == nil
}
