package rbac

type RoleStatement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
}

type Role struct {
	Version   string          `json:"Version"`
	Statement []RoleStatement `json:"Statement"`
}

func NewRoleAll(effect string) *Role {
	role := Role{Version: "2017-05-05"}
	role.Statement = make([]RoleStatement, 0)
	rs := RoleStatement{
		Effect:   effect,
		Action:   []string{"*"},
		Resource: []string{"*"},
	}

	role.Statement = append(role.Statement, rs)
	return &role
}
