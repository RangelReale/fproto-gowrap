package fproto_gowrap

import (
	"errors"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto/fdep"
)

// Default type converter for GoWrap.
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

func (tc *TypeConverter_Default) GenerateField(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	switch xfld := fld.(type) {
	case *fproto.FieldElement:
		ftype, _, scalar := g.GetType(xfld.Type, false)

		if !scalar {
			ftype = "*" + ftype
		}
		if xfld.Repeated {
			ftype = "[]" + ftype
		}

		g.Body().P(CamelCase(xfld.Name), " ", ftype)

		return true, nil
	case *fproto.MapFieldElement:
		ftype, _, scalar := g.GetType(xfld.Type, false)
		fktype, _, _ := g.GetType(xfld.KeyType, false)

		if !scalar {
			ftype = "*" + ftype
		}

		g.Body().P(CamelCase(xfld.Name), " map[", fktype, "]", ftype)

		return true, nil
	case *fproto.OneOfElement:
		return false, errors.New("OneOf is not supported")
	}
	return false, nil
}

func (tc *TypeConverter_Default) GenerateFieldImport(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	switch xfld := fld.(type) {
	case *fproto.FieldElement:
		if !xfld.Repeated {
			ftype, tp, scalar := g.GetType(xfld.Type, false)
			if !scalar && !tp.FileDep.IsSame(g.filedep) && tp.FileDep.DepType == fdep.DepType_Own {
				g.Body().P("{")
				g.Body().In()
				g.Body().P("if s.", CamelCase(xfld.Name), " != nil {")
				g.Body().In()

				g.Body().P("m.", CamelCase(xfld.Name), " = &", ftype, "{}")
				g.Body().P("err := m.", CamelCase(xfld.Name), ".Import(s.", CamelCase(xfld.Name), ")")
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
				g.Body().P("m.", CamelCase(xfld.Name), " = s.", CamelCase(xfld.Name))
			}
		} else {
			ftype, _, scalar := g.GetType(xfld.Type, false)

			g.Body().P("for _, mi := range s.", CamelCase(xfld.Name), " {")
			g.Body().In()
			if !scalar {
				g.Body().P("imp := &", ftype, "{}")
				g.Body().P("err := imp.Import(mi)")
				g.Body().P("if err != nil {")
				g.Body().In()
				g.Body().P("return err")
				g.Body().Out()
				g.Body().P("}")

				g.Body().P("m.", CamelCase(xfld.Name), " = append(m.", CamelCase(xfld.Name), ", imp)")
			} else {
				g.Body().P("m.", CamelCase(xfld.Name), " = append(m.", CamelCase(xfld.Name), ", mi)")
			}
			g.Body().Out()
			g.Body().P("}")
		}
		return true, nil
	case *fproto.MapFieldElement:
		ftype, _, scalar := g.GetType(xfld.Type, false)

		g.Body().P("for mn, mi := range s.", CamelCase(xfld.Name), " {")
		g.Body().In()
		if !scalar {
			g.Body().P("imp := &", ftype, "{}")
			g.Body().P("err := imp.Import(mi)")
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return err")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P("m.", CamelCase(xfld.Name), "[mn] = imp")
		} else {
			g.Body().P("m.", CamelCase(xfld.Name), "[mn] = mi")
		}
		g.Body().Out()
		g.Body().P("}")
	}
	return false, nil
}

func (tc *TypeConverter_Default) GenerateFieldExport(g *Generator, message *fproto.MessageElement, fld fproto.FieldElementTag) (bool, error) {
	switch xfld := fld.(type) {
	case *fproto.FieldElement:
		if !xfld.Repeated {
			_, tp, scalar := g.GetType(xfld.Type, true)
			if !scalar && !tp.FileDep.IsSame(g.filedep) && tp.FileDep.DepType == fdep.DepType_Own {
				g.Body().P("{")
				g.Body().In()
				g.Body().P("if m.", CamelCase(xfld.Name), " != nil {")
				g.Body().In()

				g.Body().P("exp, err := m.", CamelCase(xfld.Name), ".Export()")
				g.Body().P("if err != nil {")
				g.Body().In()
				g.Body().P("return nil, err")
				g.Body().Out()
				g.Body().P("}")
				g.Body().P("ret.", CamelCase(xfld.Name), " = exp")

				g.Body().Out()
				g.Body().P("}")
				g.Body().Out()
				g.Body().P("}")
			} else {
				g.Body().P("ret.", CamelCase(xfld.Name), " = m.", CamelCase(xfld.Name))
			}
		} else {
			g.Body().P("for _, mi := range m.", CamelCase(xfld.Name), " {")
			g.Body().In()
			g.Body().P("exp, err := mi.Export()")
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return nil, err")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P("ret.", CamelCase(xfld.Name), " = append(ret.", CamelCase(xfld.Name), ", exp)")
			g.Body().Out()
			g.Body().P("}")
		}
		return true, nil
	case *fproto.MapFieldElement:
		_, _, scalar := g.GetType(xfld.Type, true)

		g.Body().P("for mn, mi := range m.", CamelCase(xfld.Name), " {")
		g.Body().In()
		if !scalar {
			g.Body().P("exp, err := mi.Export()")
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return nil, err")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P("ret.", CamelCase(xfld.Name), "[mn] = exp")
		} else {
			g.Body().P("ret.", CamelCase(xfld.Name), "[mn] = mi")
		}
		g.Body().Out()
		g.Body().P("}")
	}
	return false, nil
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
