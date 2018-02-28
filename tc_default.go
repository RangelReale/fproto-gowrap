package fproto_gowrap

import (
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto/fdep"
)

type TypeConverter_Default struct {
}

func (tc *TypeConverter_Default) GetSources() []TypeConverterSource {
	return []TypeConverterSource{}
}

func (tc *TypeConverter_Default) GetType(g *Generator, fldtype string, pbsource bool) (string, bool) {
	ftype, _, scalar := g.GetType(fldtype, pbsource)
	if !scalar {
		ftype = "*" + ftype
	}
	return ftype, true
}

func (tc *TypeConverter_Default) GenerateField(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	ftype, _, scalar := g.GetType(fld.Type, false)

	if !scalar {
		ftype = "*" + ftype
	}
	if fld.Repeated {
		ftype = "[]" + ftype
	}

	g.Body().P(CamelCase(fld.Name), " ", ftype)

	return true, nil
}

func (tc *TypeConverter_Default) GenerateFieldImport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	if !fld.Repeated {
		ftype, tp, scalar := g.GetType(fld.Type, false)
		if !scalar && !tp.FileDep.IsSame(g.filedep) && tp.FileDep.DepType == fdep.DepType_Own {
			g.Body().P("{")
			g.Body().In()
			g.Body().P("if s.", CamelCase(fld.Name), " != nil {")
			g.Body().In()

			g.Body().P("m.", CamelCase(fld.Name), " = &", ftype, "{}")
			g.Body().P("err := m.", CamelCase(fld.Name), ".Import(s.", CamelCase(fld.Name), ")")
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return err")
			g.Body().Out()
			g.Body().P("}")

			g.Body().Out()
			g.Body().P("}")

			g.Body().Out()
			g.Body().P("}")
		} else {
			g.Body().P("m.", CamelCase(fld.Name), " = s.", CamelCase(fld.Name))
		}
	} else {
		ftype, _, scalar := g.GetType(fld.Type, false)

		g.Body().P("for _, mi := range s.", CamelCase(fld.Name), " {")
		g.Body().In()
		if !scalar {
			g.Body().P("imp := &", ftype, "{}")
			g.Body().P("err := imp.Import(mi)")
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return err")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P("m.", CamelCase(fld.Name), " = append(m.", CamelCase(fld.Name), ", imp)")
		} else {
			g.Body().P("m.", CamelCase(fld.Name), " = append(m.", CamelCase(fld.Name), ", mi)")
		}
		g.Body().Out()
		g.Body().P("}")
	}
	return true, nil
}

func (tc *TypeConverter_Default) GenerateFieldExport(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	if !fld.Repeated {
		_, tp, scalar := g.GetType(fld.Type, true)
		if !scalar && !tp.FileDep.IsSame(g.filedep) && tp.FileDep.DepType == fdep.DepType_Own {
			g.Body().P("{")
			g.Body().In()
			g.Body().P("if m.", CamelCase(fld.Name), " != nil {")
			g.Body().In()

			g.Body().P("exp, err := m.", CamelCase(fld.Name), ".Export()")
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return nil, err")
			g.Body().Out()
			g.Body().P("}")
			g.Body().P("ret.", CamelCase(fld.Name), " = exp")

			g.Body().Out()
			g.Body().P("}")
			g.Body().Out()
			g.Body().P("}")
		} else {
			g.Body().P("ret.", CamelCase(fld.Name), " = m.", CamelCase(fld.Name))
		}
	} else {
		g.Body().P("for _, mi := range m.", CamelCase(fld.Name), " {")
		g.Body().In()
		g.Body().P("exp, err := mi.Export()")
		g.Body().P("if err != nil {")
		g.Body().In()
		g.Body().P("return nil, err")
		g.Body().Out()
		g.Body().P("}")

		g.Body().P("ret.", CamelCase(fld.Name), " = append(ret.", CamelCase(fld.Name), ", exp)")
		g.Body().Out()
		g.Body().P("}")
	}

	return true, nil
}

func (tc *TypeConverter_Default) GenerateSrvImport(srvtype string, g *Generator, reqVarName string, retVarName string, fldtype string) (bool, error) {
	if srvtype != "grpc" {
		return false, nil
	}

	cftype_req, tp, req_scalar := g.GetType(fldtype, false)

	if tp.FileDep.DepType == fdep.DepType_Own && !req_scalar {
		g.Body().P(retVarName, " := &", cftype_req, "{}")
		g.Body().P("err = ", retVarName, ".Import(", reqVarName, ")")
	} else {
		g.Body().P(retVarName, " := ", reqVarName)
	}

	return true, nil
}

func (tc *TypeConverter_Default) GenerateSrvExport(srvtype string, g *Generator, reqVarName string, retVarName string, fldtype string) (bool, error) {
	if srvtype != "grpc" {
		return false, nil
	}

	_, tp, req_scalar := g.GetType(fldtype, false)

	if tp.FileDep.DepType == fdep.DepType_Own && !req_scalar {
		g.Body().P(retVarName, ", err := ", reqVarName, ".Export()")
	} else {
		g.Body().P(retVarName, " := ", reqVarName)
	}

	return true, nil
}

func (tc *TypeConverter_Default) EmptyValue(g *Generator, fldtype string, pbsource bool) (string, bool) {
	ftype, _, req_scalar := g.GetType(fldtype, pbsource)

	if !req_scalar {
		return "&" + ftype + "{}", true
	} else {
		return `""`, true
	}
}
