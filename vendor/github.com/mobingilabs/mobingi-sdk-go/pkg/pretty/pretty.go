package pretty

import (
	"bytes"
	"encoding/json"
)

var Pad int = 2

func Indent(count int) string {
	pad := ""
	for i := 0; i < count; i++ {
		pad += " "
	}

	return pad
}

// JSON returns a prettified JSON string of `v`.
func JSON(v interface{}, indent int) string {
	var out bytes.Buffer
	var b []byte

	pad := Indent(indent)
	_, ok := v.(string)
	if !ok {
		tmp, err := json.Marshal(v)
		if err != nil {
			return err.Error()
		}

		b = tmp
	} else {
		b = []byte(v.(string))
	}

	err := json.Indent(&out, b, "", pad)
	if err != nil {
		return err.Error()
	}

	return out.String()
}
