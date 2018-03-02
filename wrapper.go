package fproto_gowrap

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/RangelReale/fproto/fdep"
)

// Root wrapper class
type Wrapper struct {
	dep *fdep.Dep

	PkgSource PkgSource
	TypeConvs []TypeConverter
}

// Creates a new wrapper
func NewWrapper(dep *fdep.Dep) *Wrapper {
	return &Wrapper{
		dep: dep,
	}
}

// Generates one file
func (wp *Wrapper) GenerateFile(filename string, w io.Writer) error {
	g, err := NewGenerator(wp.dep, filename)
	g.PkgSource = wp.PkgSource
	g.TypeConvs = wp.TypeConvs
	if err != nil {
		return err
	}

	err = g.Generate()
	if err != nil {
		return err
	}

	fmt.Printf("*** %s\n", g.Filename())

	_, err = w.Write([]byte(g.String()))
	return err
}

// Generates all owned files.
func (wp *Wrapper) GenerateFiles(outputpath string) error {
	for _, df := range wp.dep.Files {
		if df.DepType == fdep.DepType_Own {
			g, err := NewGenerator(wp.dep, df.FilePath)
			if err != nil {
				return err
			}
			g.PkgSource = wp.PkgSource
			g.TypeConvs = wp.TypeConvs

			err = g.Generate()
			if err != nil {
				return err
			}

			p := filepath.Join(outputpath, g.Filename())

			err = os.MkdirAll(filepath.Dir(p), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(p)
			if err != nil {
				return err
			}

			_, err = file.WriteString(g.String())
			file.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
