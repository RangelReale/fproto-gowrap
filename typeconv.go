package fproto_gowrap

import "github.com/RangelReale/fproto"

// Interface for custom type converters
type TypeConverter interface {
	// Returns the list of type identifications this converter supports.
	GetSources() []TypeConverterSource
	// Returns a type name. If pbsource=true, returns the Golang protobuf type, else returns the GoWrap type.
	GetType(g *Generator, fldtype string, pbsource bool) (string, bool)
	// Generates a field inside a struct definition.
	GenerateField(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error)
	// Generates a field helper if needed (oneof needs it).
	GenerateFieldHelper(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error)
	// Generates code to import one fields.
	GenerateFieldImport(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error)
	// Generates code to export one fields.
	GenerateFieldExport(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error)

	// Generate codee to import one field into service definitions.
	GenerateSrvImport(srvtype string, g *Generator, reqVarName string, retVarName string, fldtype string) (bool, error)
	// Generate codee to export one field into service definitions.
	GenerateSrvExport(srvtype string, g *Generator, reqVarName string, retVarName string, fldtype string) (bool, error)

	// Returns the empty value for the type.
	EmptyValue(g *Generator, fldtype string, pbsource bool) (string, bool)
}

// Identification of a type
type TypeConverterSource struct {
	// protobuf file path (ex: google/protobuf/empty.proto)
	FilePath string
	// protobuf package name (ex: google.protobuf)
	PackageName string
}

// Wrapper to use more than one TypeConverter at the same time
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

func (t *TypeConverterList) GenerateField(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	for _, item := range t.list {
		ok, err := item.GenerateField(g, message, fld)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func (t *TypeConverterList) GenerateFieldHelper(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	for _, item := range t.list {
		ok, err := item.GenerateFieldHelper(g, message, fld)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func (t *TypeConverterList) GenerateFieldImport(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	for _, item := range t.list {
		ok, err := item.GenerateFieldImport(g, message, fld)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func (t *TypeConverterList) GenerateFieldExport(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
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
