package fproto_gowrap

import "github.com/RangelReale/fproto/fdep"

type PkgSource interface {
	GetPkg(filedep *fdep.FileDep) (string, bool)
}
