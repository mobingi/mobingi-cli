package sesha3

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type simplestatus struct {
	Status      string  `json:"status"`                // error, success
	Description string  `json:"description"`           // any string
	Trace       *string `json:"stack_trace,omitempty"` // any string
}

func (s simplestatus) Marshal() []byte {
	b, _ := json.Marshal(s)
	return b
}

func NewSimpleSuccess(desc string) simplestatus {
	return simplestatus{
		Status:      "success",
		Description: desc,
	}
}

func NewSimpleError(v interface{}) simplestatus {
	switch v.(type) {
	case string:
		err := fmt.Errorf("%s", v.(string))
		str := fmt.Sprintf("%+v", errors.WithStack(err))
		return simplestatus{
			Status:      "error",
			Description: v.(string),
			Trace:       &str,
		}
	case error:
		str := fmt.Sprintf("%+v", errors.WithStack(v.(error)))
		return simplestatus{
			Status:      "error",
			Description: v.(error).Error(),
			Trace:       &str,
		}
	default:
		return simplestatus{
			Status:      "error",
			Description: fmt.Sprintf("%s", v),
		}
	}
}
