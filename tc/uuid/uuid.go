package tcuuid

import (
	"fmt"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-gowrap"
)

// TypeConverter for github.com/RangelReale/go.uuid using uuid.proto.
type UUIDConvert struct {
}

func (tc *UUIDConvert) GetType(g *fproto_gowrap.Generator, fldtype string, pbsource bool) (string, bool) {
	if !pbsource {
		alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")
		return fmt.Sprintf("%s.%s", alias, "UUID"), true
	} else {
		return "", false
	}
}

func (tc *UUIDConvert) GetSources() []fproto_gowrap.TypeConverterSource {
	return []fproto_gowrap.TypeConverterSource{
		{
			FilePath:    "github.com/RangelReale/fproto-gowrap/tc/uuid/uuid.proto",
			PackageName: "uuid",
		},
	}
}

func (tc *UUIDConvert) GenerateField(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	switch xfld := fld.(type) {
	case *fproto.FieldElement:
		alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")

		ftype := fmt.Sprintf("%s.UUID", alias)

		if xfld.Repeated {
			ftype = "[]" + ftype
		}

		g.Body().P(fproto_gowrap.CamelCase(xfld.Name), " ", ftype)

		return true, nil
	}
	return false, nil
}

func (tc *UUIDConvert) GenerateFieldImport(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	switch xfld := fld.(type) {
	case *fproto.FieldElement:
		alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")

		g.Body().P("if s.", fproto_gowrap.CamelCase(xfld.Name), " != nil {")
		g.Body().In()

		g.Body().P("var err error")
		g.Body().P("m.", fproto_gowrap.CamelCase(xfld.Name), ", err = ", alias, ".FromString(s.", fproto_gowrap.CamelCase(xfld.Name), ".Value)")
		g.Body().P("if err != nil {")
		g.Body().In()
		g.Body().P("return err")
		g.Body().Out()
		g.Body().P("}")

		g.Body().Out()
		g.Body().P("}")

		return true, nil
	}
	return false, nil
}

func (tc *UUIDConvert) GenerateFieldExport(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	switch xfld := fld.(type) {
	case *fproto.FieldElement:
		ftype, _, _ := g.GetType(xfld.Type, true)

		g.Body().P("ret.", fproto_gowrap.CamelCase(xfld.Name), " = &", ftype, "{}")
		g.Body().P("ret.", fproto_gowrap.CamelCase(xfld.Name), ".Value = m.", fproto_gowrap.CamelCase(xfld.Name), ".String()")

		return true, nil
	}
	return false, nil
}

func (tc *UUIDConvert) GenerateSrvImport(srvtype string, g *fproto_gowrap.Generator, reqVarName string, retVarName string, fldtype string) (bool, error) {
	if srvtype != "grpc" {
		return false, nil
	}

	alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")

	g.Body().P("var ", retVarName, " ", alias, ".UUID")

	g.Body().P(retVarName, ", err = ", alias, ".FromString(", reqVarName, ".Value)")

	return true, nil
}

func (tc *UUIDConvert) GenerateSrvExport(srvtype string, g *fproto_gowrap.Generator, reqVarName string, retVarName string, fldtype string) (bool, error) {
	if srvtype != "grpc" {
		return false, nil
	}

	ftype, _, _ := g.GetType(fldtype, true)

	g.Body().P(retVarName, " := &", ftype, "{}")
	g.Body().P(retVarName, ".Value = ", reqVarName, ".String()")

	return true, nil
}

func (tc *UUIDConvert) EmptyValue(g *fproto_gowrap.Generator, fldtype string, pbsource bool) (string, bool) {
	if !pbsource {
		alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")
		return alias + ".UUID{}", true
	}
	return "", false
}
