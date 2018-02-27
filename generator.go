package fproto_gowrap

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto/fdep"
)

type Generator struct {
	dep     *fdep.Dep
	filedep *fdep.FileDep
	b_head  *Builder
	b_body  *Builder

	imports map[string]string
}

func NewGenerator(dep *fdep.Dep, filename string) (*Generator, error) {
	filedep, ok := dep.Files[filename]
	if !ok {
		return nil, fmt.Errorf("File %s not found", filename)
	}

	return &Generator{
		dep:     dep,
		filedep: filedep,
		b_head:  NewBuilder(),
		b_body:  NewBuilder(),
		imports: make(map[string]string),
	}, nil
}

func (g *Generator) Body() *Builder {
	return g.b_body
}

func (g *Generator) Generate() error {
	for _, message := range g.filedep.ProtoFile.Messages {
		structName := CamelCaseSlice(strings.Split(message.Name, "."))
		sourceAlias := g.Dep(g.filedep, "")

		g.b_body.P("type ", structName, " struct {")
		g.b_body.In()
		for _, fld := range message.Fields {
			ftype := g.getType(fld)

			if fld.Repeated {
				ftype = "[]" + ftype
			}

			g.b_body.P(CamelCase(fld.Name), " ", ftype)
		}
		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		// Import
		g.b_body.P("func (m *", structName, ") Import(s *", sourceAlias, ".", structName, ") error {")
		g.b_body.In()

		for _, fld := range message.Fields {
			if !fld.Repeated {
				g.b_body.P("m.", CamelCase(fld.Name), " = s.", CamelCase(fld.Name))
			} else {
				ftype := g.getType(fld)

				g.b_body.P("for _, mi := range s.", CamelCase(fld.Name), " {")
				g.b_body.In()
				g.b_body.P("var imp ", ftype)
				g.b_body.P("err := imp.Import(mi)")
				g.b_body.P("if err != nil {")
				g.b_body.In()
				g.b_body.P("return err")
				g.b_body.Out()
				g.b_body.P("}")

				g.b_body.P("m.", CamelCase(fld.Name), " = append(m.", CamelCase(fld.Name), ", imp)")
				g.b_body.Out()
				g.b_body.P("}")
			}
		}

		g.b_body.P("return nil")

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		// Export
		g.b_body.P("func (m *", structName, ") Export() (*", sourceAlias, ".", structName, ", error) {")
		g.b_body.In()

		g.b_body.P("ret := &", sourceAlias, ".", structName, "{}")

		for _, fld := range message.Fields {
			if !fld.Repeated {
				g.b_body.P("ret.", CamelCase(fld.Name), " = m.", CamelCase(fld.Name))
			} else {
				g.b_body.P("for _, mi := range m.", CamelCase(fld.Name), " {")
				g.b_body.In()
				g.b_body.P("exp, err := mi.Export()")
				g.b_body.P("if err != nil {")
				g.b_body.In()
				g.b_body.P("return nil, err")
				g.b_body.Out()
				g.b_body.P("}")

				g.b_body.P("ret.", CamelCase(fld.Name), " = append(ret.", CamelCase(fld.Name), ", exp)")
				g.b_body.Out()
				g.b_body.P("}")
			}
		}

		g.b_body.P("return ret, nil")

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

	}

	return nil
}

func (g *Generator) getType(fld *fproto.FieldElement) string {
	var ftype string

	// check if if scalar
	if st, ok := fproto.ParseScalarType(fld.Type); ok {
		ftype = st.String()
	} else {
		tp, err := g.filedep.GetType(fld.Type)
		if err != nil {
			log.Fatal(err)
		}

		if tp.FileDep.IsSame(g.filedep) {
			//_ = g.Dep(tp.FileDep, tp.Alias)
			ftype = fmt.Sprintf("*%s", tp.Name)
		} else {
			falias := g.Dep(tp.FileDep, tp.Alias)
			ftype = fmt.Sprintf("*%s.%s", falias, tp.Name)
		}
	}

	return ftype
}

func (g *Generator) Dep(filedep *fdep.FileDep, defalias string) string {
	p := filedep.GoPackage()
	var alias string
	var ok bool
	if alias, ok = g.imports[p]; ok {
		return alias
	}

	if defalias == "" {
		defalias = path.Base(p)
	}

	alias = defalias
	aliasct := 0
	aliasok := false
	for !aliasok {
		aliasok = true

		for _, a := range g.imports {
			if a == alias {
				aliasct++
				alias = fmt.Sprintf("%s%d", defalias, aliasct)
				aliasok = false
			}
		}

		if aliasok {
			break
		}
	}

	g.imports[p] = alias
	return alias
}

func (g *Generator) String() string {
	p := baseName(g.GoWrapPackage())

	g.b_head.P("package ", p)
	g.b_head.P()
	for i, ia := range g.imports {
		g.b_head.P("import ", ia, ` "`, i, `"`)
	}
	g.b_head.P()

	return g.b_head.String() + g.b_body.String()
}

func (g *Generator) Filename() string {
	p := g.GoWrapPackage()
	return path.Join(p, strings.TrimSuffix(path.Base(g.filedep.FilePath), path.Ext(g.filedep.FilePath))+".gpb.go")
}

func (g *Generator) GoWrapPackage() string {
	for _, o := range g.filedep.ProtoFile.Options {
		if o.Name == "gowrap_package" {
			return o.Value
		}
	}
	for _, o := range g.filedep.ProtoFile.Options {
		if o.Name == "go_package" {
			return o.Value
		}
	}
	return path.Dir(g.filedep.FilePath)
}
