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

func (tc *TypeConverter_Default) GenerateField(g *Generator, message *fproto.MessageElement, fld *fproto.FieldElement) (bool, error) {
	ftype, _, scalar := g.GetType(fld, false)

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
		ftype, tp, scalar := g.GetType(fld, false)
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
		ftype, _, scalar := g.GetType(fld, true)

		g.Body().P("for _, mi := range s.", CamelCase(fld.Name), " {")
		g.Body().In()
		if !scalar {
			g.Body().P("var imp *", ftype)
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
		_, tp, scalar := g.GetType(fld, true)
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
