package roles

import (
	"context"
)

type Perm uint8

func (p Perm) Ok(defaul bool) (ok bool) {
	switch p {
	case UNDEF:
		return defaul
	case ALLOW:
		return true
	default:
		return false
	}
}

func (p Perm) Allow() bool {
	return p == ALLOW
}

func (p Perm) Deny() bool {
	return p == DENY
}

func (Perm) ParseBool(v bool) Perm {
	if v {
		return ALLOW
	}
	return DENY
}

const (
	UNDEF Perm = iota
	ALLOW
	DENY
)

// Permissioner permissioner interface
type Permissioner interface {
	// HasPermissionE check has permission for permissioners or not and return error
	HasPermission(ctx context.Context, mode PermissionMode, roles ...interface{}) Perm
}

// Permissioners slice of permissioner
func Permissioners(ps ...Permissioner) Permissioner {
	return permissioners(ps)
}

type permissioners []Permissioner

func (ps permissioners) HasPermission(ctx context.Context, mode PermissionMode, roles ...interface{}) (perm Perm) {
	for _, p := range ps {
		if perm = p.HasPermission(ctx, mode, roles...); perm != UNDEF {
			return
		}
	}
	return
}