package fproto_gowrap_tc

import (
	"fmt"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-gowrap"
)

type TypeConverter_UUID struct {
}

func (tc *TypeConverter_UUID) GetType(g *fproto_gowrap.Generator, fldtype string, pbsource bool) string {
	alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")
	return fmt.Sprintf("%s.%s", alias, "UUID")
}

func (tc *TypeConverter_UUID) GetSources() []fproto_gowrap.TypeConverterSource {
	return []fproto_gowrap.TypeConverterSource{
		{
			FilePath:    "pt1/base/uuid.proto",
			PackageName: "p1_base",
		},
	}
}

func (tc *TypeConverter_UUID) GenerateField(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")

	ftype := fmt.Sprintf("%s.UUID", alias)

	if fld.Repeated {
		ftype = "[]" + ftype
	}

	g.Body().P(fproto_gowrap.CamelCase(fld.Name), " ", ftype)

	return true, nil
}

func (tc *TypeConverter_UUID) GenerateFieldImport(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
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

func (tc *TypeConverter_UUID) GenerateFieldExport(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	ftype, _, _ := g.GetType(fld.Type, true)

	g.Body().P("ret.", fproto_gowrap.CamelCase(fld.Name), " = &", ftype, "{}")
	g.Body().P("ret.", fproto_gowrap.CamelCase(fld.Name), ".Value = m.", fproto_gowrap.CamelCase(fld.Name), ".String()")

	return true, nil
}

func (tc *TypeConverter_UUID) GenerateSrvImport(g *fproto_gowrap.Generator, fldtype string) (bool, error) {
	alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")

	g.Body().P("var wreq ", alias, ".UUID")

	g.Body().P("wreq, err = ", alias, ".FromString(req.Value)")

	return true, nil
}

func (tc *TypeConverter_UUID) GenerateSrvExport(g *fproto_gowrap.Generator, fldtype string) (bool, error) {
	ftype, _, _ := g.GetType(fldtype, true)

	g.Body().P("wresp := &", ftype, "{}")
	g.Body().P("wresp.Value = resp.String()")

	return true, nil
}

func (tc *TypeConverter_UUID) EmptyValue(g *fproto_gowrap.Generator, fldtype string, pbsource bool) (string, bool) {
	if !pbsource {
		alias := g.Dep("github.com/RangelReale/go.uuid", "uuid")
		return "&" + alias + ".UUID{}", true
	}
	return "", false
}
