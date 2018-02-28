package tctime

import (
	"fmt"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-gowrap"
)

type TimeConvert struct {
}

func (tc *TimeConvert) GetSources() []fproto_gowrap.TypeConverterSource {
	return []fproto_gowrap.TypeConverterSource{
		{
			FilePath:    "google/protobuf/timestamp.proto",
			PackageName: "google.protobuf",
		},
	}
}

func (tc *TimeConvert) GetType(g *fproto_gowrap.Generator, fldtype string, pbsource bool) string {
	alias := g.Dep("time", "time")
	return fmt.Sprintf("%s.%s", alias, "Time")
}

func (tc *TimeConvert) GenerateField(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	alias := g.Dep("time", "time")

	ftype := fmt.Sprintf("%s.Time", alias)

	if fld.Repeated {
		ftype = "[]" + ftype
	}

	g.Body().P(fproto_gowrap.CamelCase(fld.Name), " ", ftype)

	return true, nil
}

func (tc *TimeConvert) GenerateFieldImport(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	pb_alias := g.Dep("github.com/golang/protobuf/ptypes", "pb_types")

	g.Body().P("if s.", fproto_gowrap.CamelCase(fld.Name), " != nil {")
	g.Body().In()

	g.Body().P("var err error")
	g.Body().P("m.", fproto_gowrap.CamelCase(fld.Name), ", err = ", pb_alias, ".Timestamp(s.", fproto_gowrap.CamelCase(fld.Name), ")")
	g.Body().P("if err != nil {")
	g.Body().In()
	g.Body().P("return err")
	g.Body().Out()
	g.Body().P("}")

	g.Body().Out()
	g.Body().P("}")

	return true, nil
}

func (tc *TimeConvert) GenerateFieldExport(g *fproto_gowrap.Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	pb_alias := g.Dep("github.com/golang/protobuf/ptypes", "pb_types")

	g.Body().P("{")
	g.Body().In()

	g.Body().P("var err error")
	g.Body().P("ret.", fproto_gowrap.CamelCase(fld.Name), ", err = ", pb_alias, ".TimestampProto(m.", fproto_gowrap.CamelCase(fld.Name), ")")
	g.Body().P("if err != nil {")
	g.Body().In()
	g.Body().P("return nil, err")
	g.Body().Out()
	g.Body().P("}")

	g.Body().Out()
	g.Body().P("}")

	return true, nil
}

func (tc *TimeConvert) GenerateSrvImport(srvtype string, g *fproto_gowrap.Generator, fldtype string) (bool, error) {
	if srvtype != "grpc" {
		return false, nil
	}

	return false, nil
}

func (tc *TimeConvert) GenerateSrvExport(srvtype string, g *fproto_gowrap.Generator, fldtype string) (bool, error) {
	if srvtype != "grpc" {
		return false, nil
	}

	return false, nil
}

func (tc *TimeConvert) EmptyValue(g *fproto_gowrap.Generator, fldtype string, pbsource bool) (string, bool) {
	if pbsource {
		alias := g.Dep("time", "time")
		return "&" + alias + ".Time{}", true
	}
	return "", false
}
