package fproto_gowrap

import "github.com/RangelReale/fproto"

type TypeConverter interface {
	GetSources() []TypeConverterSource
	GenerateField(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)
	GenerateFieldImport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)
	GenerateFieldExport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)
}

type TypeConverterSource struct {
	FilePath    string
	PackageName string
}
