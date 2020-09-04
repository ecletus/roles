package roles

import "net/http"

// Global global role instance
var Global = &Role{}

// Register register role with conditions
func Descriptors() DescriptorSlice {
	return Global.Descriptors()
}

// Register register role with conditions
func Register(descriptor ...*Descriptor) {
	Global.Register(descriptor...)
}

// Allow allows permission mode for roles
func Allow(mode PermissionMode, roles ...string) *Permission {
	return Global.Allow(mode, roles...)
}

// Deny deny permission mode for roles
func Deny(mode PermissionMode, roles ...string) *Permission {
	return Global.Deny(mode, roles...)
}

// DenyAnother deny another roles for permission mode
func DenyAnother(mode PermissionMode, roles ...string) *Permission {
	return Global.DenyAnother(mode, roles...)
}

// DenyAny deny any roles for permission mode
func AllowAny(roles ...string) *Permission {
	return Global.AllowAny(roles...)
}

// Get role defination
func Get(name string) (*Descriptor, bool) {
	return Global.Get(name)
}

// MustGet role defination
func MustGet(name string) (d *Descriptor) {
	d, _ = Global.Get(name)
	return
}

// Remove role definition from global role instance
func Remove(name string) {
	Global.Remove(name)
}

// Reset role definitions from global role instance
func Reset() {
	Global.Reset()
}

// MatchedRoles return defined roles from user
func MatchedRoles(req *http.Request, user interface{}) Roles {
	return Global.MatchedRoles(req, user)
}

// HasRole check if current user has role
func HasRole(req *http.Request, user interface{}, roles ...string) bool {
	return Global.HasRole(req, user, roles...)
}

// NewPermission initialize a new permission for default role
func NewPermission() *Permission {
	return Global.NewPermission()
}
