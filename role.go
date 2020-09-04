package roles

import (
	"fmt"
	"net/http"
)

const (
	// Anyone is a role for any one
	Anyone  = "*"
	Visitor = "visitor"
)

func GetVisitor() *Descriptor {
	return &Descriptor{
		Name: Visitor,
		Checker: func(req *http.Request, user interface{}) bool {
			return user == nil
		},
	}
}

// Checker check current request match this role or not
type Checker func(req *http.Request, user interface{}) bool

// New initialize a new `Role`
func New() *Role {
	return &Role{}
}

// Role is a struct contains all roles descriptors
type Role struct {
	descriptors map[string]*Descriptor
}

// Register register role with conditions
func (role *Role) Register(descriptor ...*Descriptor) {
	for _, descriptor := range descriptor {
		if descriptor.Checker == nil {
			panic("checker is nil")
		}
		if role.descriptors == nil {
			role.descriptors = map[string]*Descriptor{}
		}
		if role.descriptors[descriptor.Name] != nil {
			fmt.Printf("%v already defined, overwrited it!\n", descriptor.Name)
		}
		role.descriptors[descriptor.Name] = descriptor
	}
}

// DescriptorSlice return slice of descriptors
func (role *Role) Descriptors() (desciptors DescriptorSlice) {
	if role.descriptors != nil {
		for _, descriptor := range role.descriptors {
			desciptors = append(desciptors, descriptor)
		}
	}
	return
}

// NewPermission initialize permission
func (role *Role) NewPermission() *Permission {
	return &Permission{
		Role:               role,
		AllowedAny:         map[string]bool{},
		AllowedRoles:       map[PermissionMode][]string{},
		DeniedRoles:        map[PermissionMode][]string{},
		DaniedAnotherRoles: map[PermissionMode][]string{},
	}
}

// Allow allows permission mode for roles
func (role *Role) Allow(mode PermissionMode, roles ...string) *Permission {
	return role.NewPermission().Allow(mode, roles...)
}

// Deny deny permission mode for roles
func (role *Role) Deny(mode PermissionMode, roles ...string) *Permission {
	return role.NewPermission().Deny(mode, roles...)
}

// DenyAnother deny another roles for permission mode
func (role *Role) DenyAnother(mode PermissionMode, roles ...string) *Permission {
	return role.NewPermission().DenyAnother(mode, roles...)
}

// DenyAny deny all roles for permission mode
func (role *Role) AllowAny(roles ...string) *Permission {
	return role.NewPermission().AllowAny(roles...)
}

// Get role defination
func (role *Role) Get(name string) (descriptor *Descriptor, ok bool) {
	descriptor, ok = role.descriptors[name]
	return
}

// Roles return roles names
func (role *Role) Roles() (roles Roles) {
	roles.m = map[string]interface{}{}
	for name := range role.descriptors {
		roles.Append(name)
	}
	return
}

// Remove role descriptor
func (role *Role) Remove(name string) {
	delete(role.descriptors, name)
}

// Reset role descriptors
func (role *Role) Reset() {
	role.descriptors = map[string]*Descriptor{}
}

// MatchedRoles return defined roles from user
func (role *Role) MatchedRoles(req *http.Request, user interface{}) (roles Roles) {
	if descriptors := role.descriptors; descriptors != nil {
		for name, descriptor := range descriptors {
			if descriptor.Checker(req, user) {
				roles.Append(name)
			}
		}
	}
	return
}

// HasRole check if current user has role
func (role *Role) HasRole(req *http.Request, user interface{}, roles ...string) bool {
	if descriptors := role.descriptors; descriptors != nil {
		for _, name := range roles {
			if descriptor, ok := descriptors[name]; ok {
				if descriptor.Checker(req, user) {
					return true
				}
			}
		}
	}
	return false
}

// Copy Copy this role
func (role Role) Copy() *Role {
	var descriptors = map[string]*Descriptor{}
	for k, v := range role.descriptors {
		descriptors[k] = v
	}
	role.descriptors = descriptors
	return &role
}
