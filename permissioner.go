package roles

import "errors"

var ErrDefaultPermission = errors.New("default permission")

func IsDefaultPermission(err error) bool {
	if err == nil {
		return false
	}
	return err == ErrDefaultPermission
}

// Permissioner permissioner interface
type Permissioner interface {
	// HasPermissionE check has permission for permissioners or not and return error
	HasPermissionE(mode PermissionMode, roles ...interface{}) (ok bool, err error)
}

// ConcatPermissioner concat permissioner
func ConcatPermissioner(ps ...Permissioner) Permissioner {
	var newPS []Permissioner
	for _, p := range ps {
		if p != nil {
			newPS = append(newPS, p)
		}
	}
	return permissioners(newPS)
}

type permissioners []Permissioner

func (ps permissioners) HasPermissionE(mode PermissionMode, roles ...interface{}) (ok bool, err error) {
	for _, p := range ps {
		if p != nil && !HasPermission(p, mode, roles...) {
			return
		}
	}

	return true, ErrDefaultPermission
}

func HasPermission(permissioner Permissioner, mode PermissionMode, roles ...interface{}) (ok bool) {
	ok, _ = permissioner.HasPermissionE(mode, roles...)
	return
}

func HasPermissionDefault(defaul bool, permissioner Permissioner, mode PermissionMode, roles ...interface{}) (ok bool) {
	ok, _ = HasPermissionDefaultE(defaul, permissioner, mode, roles...)
	return
}

func HasPermissionDefaultE(defaul bool, permissioner Permissioner, mode PermissionMode, roles ...interface{}) (ok bool, err error) {
	ok, err = permissioner.HasPermissionE(mode, roles...)
	if IsDefaultPermission(err) {
		ok = defaul
	}
	return
}
