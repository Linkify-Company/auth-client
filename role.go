package domain

type Role int

const (
	RoleAdmin = iota
	RoleModerator
	RoleUser
)

func (r *Role) SetDefault() {
	*r = RoleUser
}
