package fproto_gowrap

import "github.com/RangelReale/fproto"

type TypeConverter interface {
	GetSources() []TypeConverterSource
	GetType(g *Generator, fldtype string, pbsource bool) string
	GenerateField(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)
	GenerateFieldImport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)
	GenerateFieldExport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)

	GenerateSrvImport(g *Generator, fldtype string) (bool, error)
	GenerateSrvExport(g *Generator, fldtype string) (bool, error)

	EmptyValue(g *Generator, fldtype string, pbsource bool) (string, bool)
}

type TypeConverterSource struct {
	FilePath    string
	PackageName string
}
