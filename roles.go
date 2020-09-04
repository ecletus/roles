package roles

import (
	"strings"

	"github.com/moisespsena-go/i18n-modular/i18nmod"
)

type Roles struct {
	m map[string]interface{}
}

func NewRoles(roles ...string) (r Roles) {
	r.m = map[string]interface{}{}
	for _, role := range roles {
		r.m[role] = nil
	}
	return r
}

func (this *Roles) Merge(roles Roles) {
	for name := range roles.m {
		this.m[name] = nil
	}
}

func (this *Roles) Append(name ...string) {
	if this.m == nil {
		this.m = map[string]interface{}{}
	}
	for _, name := range name {
		this.m[name] = nil
	}
}

func (this Roles) Each(cb func(mode string)) {
	for mode := range this.m {
		cb(mode)
	}
}
func (this Roles) Filter(cb func(mode string) (include bool)) (roles Roles) {
	roles.m = map[string]interface{}{}
	for mode := range this.m {
		if cb(mode) {
			roles.m[mode] = nil
		}
	}
	return
}

func (this Roles) Intersection(names []string) (roles Roles) {
	roles.m = map[string]interface{}{}
	for _, name := range names {
		if _, ok := this.m[name]; ok {
			roles.m[name] = nil
		}
	}
	return
}

func (this Roles) Copy() (roles Roles) {
	roles.m = map[string]interface{}{}
	for name := range this.m {
		roles.m[name] = nil
	}
	return
}

func (this Roles) Interfaces() (result []interface{}) {
	result = make([]interface{}, len(this.m), len(this.m))
	i := 0
	for role := range this.m {
		result[i] = role
		i++
	}
	return
}

func (this Roles) Strings() (result []string) {
	result = make([]string, len(this.m), len(this.m))
	i := 0
	for role := range this.m {
		result[i] = role
		i++
	}
	return
}
func (this Roles) Has(name string) (ok bool) {
	_, ok = this.m[name]
	return
}

func (this Roles) String() string {
	return strings.Join(this.Strings(),", ")
}

func (this Roles) Len() int {
	return len(this.m)
}

func (this Roles) Descriptors() (descriptors DescriptorSlice) {
	for name := range this.m {
		if d, ok := Global.descriptors[name]; ok {
			descriptors = append(descriptors, d)
		}
	}
	return
}

func (this Roles) Labels(ctx i18nmod.Context) (labels []string) {
	for name := range this.m {
		if d, ok := Global.descriptors[name]; ok {
			labels = append(labels, d.Translate(ctx))
		}
	}
	return
}