package tcuuid

import (
	"fmt"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-gowrap"
)

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

func (tc *UUIDConvert) GenerateField(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")

	ftype := fmt.Sprintf("%s.UUID", alias)

	if fld.Repeated {
		ftype = "[]" + ftype
	}

	g.Body().P(fproto_gowrap.CamelCase(fld.Name), " ", ftype)

	return true, nil
}

func (tc *UUIDConvert) GenerateFieldImport(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")

	g.Body().P("if s.", fproto_gowrap.CamelCase(fld.Name), " != nil {")
	g.Body().In()

	g.Body().P("var err error")
	g.Body().P("m.", fproto_gowrap.CamelCase(fld.Name), ", err = ", alias, ".FromString(s.", fproto_gowrap.CamelCase(fld.Name), ".Value)")
	g.Body().P("if err != nil {")
	g.Body().In()
	g.Body().P("return err")
	g.Body().Out()
	g.Body().P("}")

	g.Body().Out()
	g.Body().P("}")

	return true, nil
}

func (tc *UUIDConvert) GenerateFieldExport(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	ftype, _, _ := g.GetType(fld.Type, true)

	g.Body().P("ret.", fproto_gowrap.CamelCase(fld.Name), " = &", ftype, "{}")
	g.Body().P("ret.", fproto_gowrap.CamelCase(fld.Name), ".Value = m.", fproto_gowrap.CamelCase(fld.Name), ".String()")

	return true, nil
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
