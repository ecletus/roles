package roles

import (
	"github.com/moisespsena-go/i18n-modular/i18nmod"

	"github.com/moisespsena-go/aorm"
)

type DescriptorSlice []*Descriptor

type Descriptor struct {
	Name  string
	Label string
	Checker Checker
}

func NewDescriptor(name string, checker Checker, label ...string) *Descriptor {
	if len(label) == 0 {
		label = []string{""}
	}
	return &Descriptor{Name: name, Checker: checker, Label: label[0]}
}

func (this Descriptor) GetID() aorm.ID {
	return aorm.NewValuedId(aorm.StrId(this.Name))
}

func (this Descriptor) Translate(ctx i18nmod.Context) string {
	if this.Label == "" {
		return this.Name
	}
	return ctx.T(this.Label).Get()
}

func (this DescriptorSlice) Each(cb func(descriptor *Descriptor)) {
	for _, descriptor := range this {
		cb(descriptor)
	}
}
func (this DescriptorSlice) Filter(cb func(descriptor *Descriptor) (include bool)) (result DescriptorSlice) {
	for _, descriptor := range this {
		if cb(descriptor) {
			result = append(result, descriptor)
		}
	}
	return
}

func (this DescriptorSlice) Names() (result []string) {
	result = make([]string, len(this), len(this))
	for i, descriptor := range this {
		result[i] = descriptor.Name
	}
	return
}

func (this DescriptorSlice) Intersection(names []string) DescriptorSlice {
	return this.Filter(func(descriptor *Descriptor) (include bool) {
		for _, name := range names {
			if name == descriptor.Name {
				return true
			}
		}
		return false
	})
}
