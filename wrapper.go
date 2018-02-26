package fproto_gowrap

import (
	"github.com/RangelReale/fproto/fdep"
)

type Wrapper struct {
	dep *fdep.Dep
}

func NewWrapper(dep *fdep.Dep) *Wrapper {
	return &Wrapper{
		dep: dep,
	}
}
