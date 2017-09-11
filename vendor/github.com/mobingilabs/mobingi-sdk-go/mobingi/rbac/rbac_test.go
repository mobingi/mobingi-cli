package rbac

import (
	"os"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestCreateRoleDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
		})

		rbac := New(sess)
		in := &CreateRoleInput{
			Name:  "testrole",
			Scope: *(NewRoleAll("Allow")),
		}

		resp, body, err := rbac.CreateRole(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_, _ = resp, body
	}
}

func TestDescribeRolesDevAcct(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
		})

		rbac := New(sess)
		in := &DescribeRolesInput{
			User: "chewsubuser1",
		}

		resp, body, err := rbac.DescribeRoles(nil)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		resp, body, err = rbac.DescribeRoles(in)
		if err != nil {
			t.Errorf("Expecting nil error, received %v", err)
		}

		_, _ = resp, body
	}
}
