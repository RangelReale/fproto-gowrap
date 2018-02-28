package fproto_gowrap

import "github.com/RangelReale/fproto"

type TypeConverter interface {
	GetSources() []TypeConverterSource
	GetType(g *Generator, fldtype string, pbsource bool) (string, bool)
	GenerateField(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)
	GenerateFieldImport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)
	GenerateFieldExport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error)

	GenerateSrvImport(srvtype string, g *Generator, reqVarName string, retVarName string, fldtype string) (bool, error)
	GenerateSrvExport(srvtype string, g *Generator, reqVarName string, retVarName string, fldtype string) (bool, error)

	EmptyValue(g *Generator, fldtype string, pbsource bool) (string, bool)
}

type TypeConverterSource struct {
	FilePath    string
	PackageName string
}

type TypeConverterList struct {
	list []TypeConverter
}

func NewTypeConverterList(list []TypeConverter) *TypeConverterList {
	return &TypeConverterList{list}
}

func (t *TypeConverterList) GetSources() []TypeConverterSource {
	if len(t.list) > 0 {
		return t.list[0].GetSources()
	}
	return []TypeConverterSource{}
}

func (t *TypeConverterList) GetType(g *Generator, fldtype string, pbsource bool) (string, bool) {
	for _, item := range t.list {
		s, ok := item.GetType(g, fldtype, pbsource)
		if ok {
			return s, ok
		}
	}
	return "", false
}

func (t *TypeConverterList) GenerateField(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	for _, item := range t.list {
		ok, err := item.GenerateField(g, message, fld)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func (t *TypeConverterList) GenerateFieldImport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	for _, item := range t.list {
		ok, err := item.GenerateFieldImport(g, message, fld)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func (t *TypeConverterList) GenerateFieldExport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	for _, item := range t.list {
		ok, err := item.GenerateFieldExport(g, message, fld)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil

}

func (t *TypeConverterList) GenerateSrvImport(srvtype string, g *Generator, reqVarName string, retVarName string, fldtype string) (bool, error) {
	for _, item := range t.list {
		ok, err := item.GenerateSrvImport(srvtype, g, reqVarName, retVarName, fldtype)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func (t *TypeConverterList) GenerateSrvExport(srvtype string, g *Generator, reqVarName string, retVarName string, fldtype string) (bool, error) {
	for _, item := range t.list {
		ok, err := item.GenerateSrvExport(srvtype, g, reqVarName, retVarName, fldtype)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func (t *TypeConverterList) EmptyValue(g *Generator, fldtype string, pbsource bool) (string, bool) {
	for _, item := range t.list {
		s, ok := item.EmptyValue(g, fldtype, pbsource)
		if ok {
			return s, ok
		}
	}
	return "", false

}
