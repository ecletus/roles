package roles

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// PermissionMode permission mode
type PermissionMode string

const (
	// Create predefined permission mode, create permission
	Create PermissionMode = "create"
	// Read predefined permission mode, read permission
	Read PermissionMode = "read"
	// Update predefined permission mode, update permission
	Update PermissionMode = "update"
	// Delete predefined permission mode, deleted permission
	Delete PermissionMode = "delete"
	// CRUD predefined permission mode, create+read+update+delete permission
	CRUD PermissionMode = "crud"
	NONE PermissionMode = ""
)

// ErrPermissionDenied no permission error
var ErrPermissionDenied = errors.New("permission denied")

// Permission a struct contains permission definitions
type Permission struct {
	Role               *Role
	AllowedAny         map[string]bool
	AllowedRoles       map[PermissionMode][]string
	DeniedRoles        map[PermissionMode][]string
	DaniedAnotherRoles map[PermissionMode][]string
}

func includeRoles(roles []string, values []string) bool {
	for _, role := range roles {
		if role == Anyone {
			return true
		}

		for _, value := range values {
			if value == role {
				return true
			}
		}
	}
	return false
}

// Concat concat two permissions into a new one
func (permission *Permission) Concat(newPermission *Permission) *Permission {
	var result = Permission{
		Role:               Global,
		AllowedAny:         map[string]bool{},
		AllowedRoles:       map[PermissionMode][]string{},
		DeniedRoles:        map[PermissionMode][]string{},
		DaniedAnotherRoles: map[PermissionMode][]string{},
	}

	var appendRoles = func(p *Permission) {
		if p != nil {
			result.Role = p.Role

			for mode, roles := range p.AllowedRoles {
				result.AllowedRoles[mode] = append(result.AllowedRoles[mode], roles...)
			}

			for mode, roles := range p.DeniedRoles {
				result.DeniedRoles[mode] = append(result.DeniedRoles[mode], roles...)
			}

			for mode, roles := range p.DaniedAnotherRoles {
				result.DaniedAnotherRoles[mode] = append(result.DaniedAnotherRoles[mode], roles...)
			}

			for role, v := range p.AllowedAny {
				result.AllowedAny[role] = v
			}
		}
	}

	appendRoles(newPermission)
	appendRoles(permission)
	return &result
}

// Allow allows permission mode for roles
func (permission *Permission) Allow(mode PermissionMode, roles ...string) *Permission {
	if mode == CRUD {
		return permission.Allow(Create, roles...).Allow(Update, roles...).Allow(Read, roles...).Allow(Delete, roles...)
	}

	if permission.AllowedRoles[mode] == nil {
		permission.AllowedRoles[mode] = []string{}
	}
	permission.AllowedRoles[mode] = append(permission.AllowedRoles[mode], roles...)
	return permission
}

// Deny deny permission mode for roles
func (permission *Permission) Deny(mode PermissionMode, roles ...string) *Permission {
	if mode == CRUD {
		return permission.Deny(Create, roles...).Deny(Update, roles...).Deny(Read, roles...).Deny(Delete, roles...)
	}

	if permission.DeniedRoles[mode] == nil {
		permission.DeniedRoles[mode] = []string{}
	}
	permission.DeniedRoles[mode] = append(permission.DeniedRoles[mode], roles...)
	return permission
}

// DenyAnother deny another roles for permission mode
func (permission *Permission) DenyAnother(mode PermissionMode, roles ...string) *Permission {
	if mode == CRUD {
		return permission.Allow(Create, roles...).Allow(Update, roles...).Allow(Read, roles...).Allow(Delete, roles...)
	}

	if permission.DaniedAnotherRoles[mode] == nil {
		permission.DaniedAnotherRoles[mode] = []string{}
	}
	permission.DaniedAnotherRoles[mode] = append(permission.DaniedAnotherRoles[mode], roles...)
	return permission
}

// DenyAnother deny any roles for permission mode
func (permission *Permission) AllowAny(roles ...string) *Permission {
	for _, role := range roles {
		permission.AllowedAny[role] = true
	}
	return permission
}

// HasPermissionS check roles strings has permission for mode or not
func (permission Permission) HasPermissionS(ctx context.Context, mode PermissionMode, roles ...string) Perm {
	rolesi := make([]interface{}, len(roles))
	for i, v := range roles {
		rolesi[i] = v
	}
	return permission.HasPermission(ctx, mode, rolesi...)
}

func (permission Permission) HasPermission(ctx context.Context, mode PermissionMode, roles ...interface{}) (perm Perm) {
	var roleNames []string
	for _, role := range roles {
		switch t := role.(type) {
		case string:
			roleNames = append(roleNames, t)
		case []string:
			roleNames = append(roleNames, t...)
		case Roler:
			roleNames = append(roleNames, t.GetRoles()...)
		default:
			panic(errors.New(fmt.Sprintf("invalid role %#v", role)))
		}
	}

	if len(permission.AllowedAny) > 0 {
		var ok2 bool
		for _, role := range roleNames {
			if ok2 = permission.AllowedAny[role]; ok2 {
				break
			}
		}
		if !ok2 {
			return DENY
		}
	}

	if len(permission.DaniedAnotherRoles) != 0 {
		if roles := permission.DaniedAnotherRoles[mode]; len(roles) > 0 {
			if !includeRoles(roles, roleNames) {
				return DENY
			}
		}
	}

	if len(permission.DeniedRoles) != 0 {
		if DeniedRoles := permission.DeniedRoles[mode]; len(DeniedRoles) > 0 {
			if includeRoles(DeniedRoles, roleNames) {
				return DENY
			}
		}
	}

	// return true if haven't define allowed roles
	if len(permission.AllowedRoles) == 0 {
		return ALLOW
	}

	if AllowedRoles := permission.AllowedRoles[mode]; len(AllowedRoles) > 0 {
		if includeRoles(AllowedRoles, roleNames) {
			return ALLOW
		}
	}
	return
}
