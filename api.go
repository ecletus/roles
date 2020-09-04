package roles

type RoleHaver interface {
	HasRole(name string) bool
}
