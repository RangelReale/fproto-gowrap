package fproto_gowrap

import (
	"io"

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

func (wp *Wrapper) GenerateFile(filename string, w io.Writer) error {
	g, err := NewGenerator(wp.dep, filename)
	if err != nil {
		return err
	}

	err = g.Generate()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(g.String()))
	return err
}
