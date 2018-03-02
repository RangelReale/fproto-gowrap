package fproto_gowrap

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto/fdep"
)

// Generators generates a wrapper for a single file.
type Generator struct {
	dep        *fdep.Dep
	filedep    *fdep.FileDep
	b_head     *Builder
	b_body     *Builder
	tc_default TypeConverter

	imports map[string]string

	// Interface to do package name generation.
	PkgSource PkgSource

	// List of type conversions
	TypeConvs []TypeConverter

	// Service generation type (default is "grpc")
	SrvType string
}

// Creates a new generator for the file path.
func NewGenerator(dep *fdep.Dep, filepath string) (*Generator, error) {
	filedep, ok := dep.Files[filepath]
	if !ok {
		return nil, fmt.Errorf("File %s not found", filepath)
	}

	return &Generator{
		dep:        dep,
		filedep:    filedep,
		b_head:     NewBuilder(),
		b_body:     NewBuilder(),
		imports:    make(map[string]string),
		tc_default: &TypeConverter_Default{},
		SrvType:    "grpc",
	}, nil
}

// Returns the body part builder
func (g *Generator) Body() *Builder {
	return g.b_body
}

// Executes the generator
func (g *Generator) Generate() error {
	err := g.GenerateMessages()
	if err != nil {
		return err
	}

	err = g.GenerateServices()
	if err != nil {
		return err
	}

	return nil
}

// Generates the protobuf messages
func (g *Generator) GenerateMessages() error {
	for _, message := range g.filedep.ProtoFile.Messages {
		structName := CamelCaseSlice(strings.Split(message.Name, "."))
		sourceAlias := g.FileDep(g.filedep, "", true)

		//
		// type MyMessage struct
		//
		g.b_body.P("type ", structName, " struct {")
		g.b_body.In()
		for _, fld := range message.Fields {
			tc := g.getTypeConv(fld.Type)
			if tc != nil {
				var err error
				_, err = tc.GenerateField(g, message, fld)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("No type converter found")
			}
		}
		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		//
		// func (m *MyMessage) Import(s *myapp.MyMessage) error
		//

		g.b_body.P("func (m *", structName, ") Import(s *", sourceAlias, ".", structName, ") error {")
		g.b_body.In()

		for _, fld := range message.Fields {
			tc := g.getTypeConv(fld.Type)
			if tc != nil {
				var err error
				_, err = tc.GenerateFieldImport(g, message, fld)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("No type converter found")
			}
		}

		g.b_body.P("return nil")

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		//
		// func (m *MyMessage) Export() (*myapp.MyMessage, error)
		//

		g.b_body.P("func (m *", structName, ") Export() (*", sourceAlias, ".", structName, ", error) {")
		g.b_body.In()

		g.b_body.P("ret := &", sourceAlias, ".", structName, "{}")

		for _, fld := range message.Fields {
			tc := g.getTypeConv(fld.Type)
			if tc != nil {
				var err error
				_, err = tc.GenerateFieldExport(g, message, fld)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("No type converter found")
			}
		}

		g.b_body.P("return ret, nil")

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

	}

	return nil
}

// Generate gRPC service wrappers.
func (g *Generator) GenerateServices() error {
	if g.SrvType == "" {
		return nil
	}
	if g.SrvType != "grpc" {
		return fmt.Errorf("Unknown service type '%s'", g.SrvType)
	}

	if len(g.filedep.ProtoFile.Services) == 0 {
		return nil
	}

	// import all required dependencies
	ctx_alias := g.Dep("golang.org/x/net/context", "context")
	grpc_alias := g.Dep("google.golang.org/grpc", "grpc")
	wraputil_alias := g.Dep("github.com/RangelReale/fproto-gowrap/wraputil", "wraputil")
	func_alias := g.FileDep(g.filedep, "", true)

	for _, service := range g.filedep.ProtoFile.Services {
		svcName := CamelCase(service.Name)

		//
		// CLIENT
		//

		//
		// type MyServiceClient interface
		//

		g.b_body.P("type ", svcName, "Client interface {")
		g.b_body.In()

		for _, rpc := range service.RPCs {
			tc_req := g.getTypeConv(rpc.RequestType)
			tc_resp := g.getTypeConv(rpc.ResponseType)

			if tc_req == nil || tc_resp == nil {
				return fmt.Errorf("No type converter found")
			}

			ftype_req, _ := tc_req.GetType(g, rpc.RequestType, false)
			ftype_resp, _ := tc_resp.GetType(g, rpc.ResponseType, false)

			//
			// MyRPC(ctx context.Context, in *MyReq, opts ...grpc.CallOption) (*MyResp, error)
			//

			g.b_body.P(rpc.Name, "(ctx ", ctx_alias, ".Context, in ", ftype_req, ", opts ...", grpc_alias, ".CallOption) (", ftype_resp, ", error)")
		}

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		//
		// type wrapMyServiceClient struct
		//

		g.b_body.P("type wrap", svcName, "Client struct {")
		g.b_body.In()

		// the default Golang protobuf client
		g.b_body.P("cli ", func_alias, ".", svcName, "Client")
		// the customizable error handler
		g.b_body.P("errorHandler ", wraputil_alias, ".ServiceErrorHandler")

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		//
		// func NewMyServiceClient(cc *grpc.ClientConn, errorHandler ...wraputil.ServiceErrorHandler) MyServiceClient
		//

		g.b_body.P("func New", svcName, "Client(cc *", grpc_alias, ".ClientConn, errorHandler ...", wraputil_alias, ".ServiceErrorHandler) ", svcName, "Client {")
		g.b_body.In()

		g.b_body.P("w := &wrap", svcName, "Client{cli: ", func_alias, ".New", svcName, "Client(cc)}")

		// check if implements ServiceErrorHandler
		g.Body().P("if len(errorHandler) > 0 {")
		g.b_body.In()
		g.Body().P("w.errorHandler = errorHandler[0]")
		g.b_body.Out()
		g.Body().P("} else {")
		g.b_body.In()
		g.Body().P("w.errorHandler = &", wraputil_alias, ".ServiceErrorHandler_Default{}")
		g.b_body.Out()
		g.Body().P("}")

		g.b_body.P("return w")

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		// Implement each RPC wrapper

		for _, rpc := range service.RPCs {
			tc_req := g.getTypeConv(rpc.RequestType)
			tc_resp := g.getTypeConv(rpc.ResponseType)

			if tc_req == nil || tc_resp == nil {
				return fmt.Errorf("No type converter found")
			}

			ftype_req, _ := tc_req.GetType(g, rpc.RequestType, false)
			ftype_resp, _ := tc_resp.GetType(g, rpc.ResponseType, false)

			//
			// func (w *wrapMyServiceClient) MyRPC(ctx context.Context, in *MyReq, opts ...grpc.CallOption) (*MyResp, error)
			//

			g.b_body.P("func (w *wrap", svcName, "Client) ", rpc.Name, "(ctx ", ctx_alias, ".Context, in ", ftype_req, ", opts ...", grpc_alias, ".CallOption) (", ftype_resp, ", error) {")
			g.b_body.In()
			g.Body().P("var err error")

			g.Body().P()

			// default return value
			defretvalue, _ := tc_resp.EmptyValue(g, rpc.ResponseType, false)

			// convert request
			_, err := tc_req.GenerateSrvExport(g.SrvType, g, "in", "wreq", rpc.RequestType)
			if err != nil {
				return err
			}

			// check error
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return ", defretvalue, ", w.errorHandler.HandleServiceError(", wraputil_alias, ".SET_EXPORT, err)")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P()

			// call
			g.Body().P("resp, err := w.cli.", rpc.Name, "(ctx, wreq)")
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return ", defretvalue, ", w.errorHandler.HandleServiceError(", wraputil_alias, ".SET_CALL, err)")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P()

			// convert response
			g.Body().P("err = nil")

			_, err = tc_resp.GenerateSrvImport(g.SrvType, g, "resp", "wresp", rpc.ResponseType)
			if err != nil {
				return err
			}

			// check error
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return ", defretvalue, ", w.errorHandler.HandleServiceError(", wraputil_alias, ".SET_IMPORT, err)")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P()

			// Return response
			g.b_body.P("return wresp, nil")

			g.b_body.Out()
			g.b_body.P("}")
			g.b_body.P()
		}

		//
		// SERVER
		//

		//
		// type MyServiceServer interface
		//

		g.b_body.P("type ", svcName, "Server interface {")
		g.b_body.In()

		for _, rpc := range service.RPCs {
			tc_req := g.getTypeConv(rpc.RequestType)
			tc_resp := g.getTypeConv(rpc.ResponseType)

			if tc_req == nil || tc_resp == nil {
				return fmt.Errorf("No type converter found")
			}

			ftype_req, _ := tc_req.GetType(g, rpc.RequestType, false)
			ftype_resp, _ := tc_resp.GetType(g, rpc.ResponseType, false)

			//
			// MyRPC(ctx.Context, *MyReq) (*MyResp, error)
			//

			g.b_body.P(rpc.Name, "(", ctx_alias, ".Context, ", ftype_req, ") (", ftype_resp, ", error)")
		}

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		//
		// type wrapMyServiceServer struct
		//

		g.b_body.P("type wrap", svcName, "Server struct {")
		g.b_body.In()

		g.b_body.P("srv ", svcName, "Server")
		g.b_body.P("errorHandler ", wraputil_alias, ".ServiceErrorHandler")

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		//
		// func newWrapMyServiceServer(srv MyServiceServer) *wrapMyServiceServer
		//

		g.b_body.P("func newWrap", svcName, "Server(srv ", svcName, "Server) *wrap", svcName, "Server {")
		g.b_body.In()

		g.b_body.P("w := &wrap", svcName, "Server{srv: srv}")

		// check if implements ServiceErrorHandler
		g.Body().P("if eh, ok := srv.(", wraputil_alias, ".ServiceErrorHandler); ok {")
		g.b_body.In()
		g.Body().P("w.errorHandler = eh")
		g.b_body.Out()
		g.Body().P("} else {")
		g.b_body.In()
		g.Body().P("w.errorHandler = &", wraputil_alias, ".ServiceErrorHandler_Default{}")
		g.b_body.Out()
		g.Body().P("}")

		g.b_body.P("return w")

		g.b_body.Out()
		g.b_body.P("}")
		g.b_body.P()

		// Generate RPCs
		for _, rpc := range service.RPCs {
			tc_req := g.getTypeConv(rpc.RequestType)
			tc_resp := g.getTypeConv(rpc.ResponseType)

			if tc_req == nil || tc_resp == nil {
				return fmt.Errorf("No type converter found")
			}

			ftype_req, _ := tc_req.GetType(g, rpc.RequestType, true)
			ftype_resp, _ := tc_resp.GetType(g, rpc.ResponseType, true)

			//
			// func (w *wrapMyServiceServer) MyRPC(ctx context.Context, req *myapp.MyReq) (*myapp.MyResp, error)
			//

			g.b_body.P("func (w *wrap", svcName, "Server) ", rpc.Name, "(ctx ", ctx_alias, ".Context, req ", ftype_req, ") (", ftype_resp, ", error) {")
			g.b_body.In()
			g.Body().P("var err error")

			g.Body().P()

			// default return value
			defretvalue, _ := tc_resp.EmptyValue(g, rpc.ResponseType, true)

			// convert request
			_, err := tc_req.GenerateSrvImport(g.SrvType, g, "req", "wreq", rpc.RequestType)
			if err != nil {
				return err
			}

			// check error
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return ", defretvalue, ", w.errorHandler.HandleServiceError(", wraputil_alias, ".SET_IMPORT, err)")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P()

			// call
			g.Body().P("resp, err := w.srv.", rpc.Name, "(ctx, wreq)")
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return ", defretvalue, ", w.errorHandler.HandleServiceError(", wraputil_alias, ".SET_CALL, err)")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P()

			// convert response
			g.Body().P("err = nil")

			_, err = tc_resp.GenerateSrvExport(g.SrvType, g, "resp", "wresp", rpc.ResponseType)
			if err != nil {
				return err
			}

			// check error
			g.Body().P("if err != nil {")
			g.Body().In()
			g.Body().P("return ", defretvalue, ", w.errorHandler.HandleServiceError(", wraputil_alias, ".SET_EXPORT, err)")
			g.Body().Out()
			g.Body().P("}")

			g.Body().P()

			// return response
			g.b_body.P("return wresp, nil")

			g.b_body.Out()
			g.b_body.P("}")
			g.b_body.P()
		}

		//
		// func RegisterMyServiceServer(s *grpc.Server, srv MyServiceServer)
		//

		g.b_body.P("func Register", svcName, "Server(s *", grpc_alias, ".Server, srv ", svcName, "Server) {")
		g.b_body.In()

		// myapp.RegisterMyServiceServer(s, newWrapMyServiceServer(srv))
		g.b_body.P(func_alias, ".Register", svcName, "Server(s, newWrap", svcName, "Server(srv))")

		g.b_body.Out()
		g.b_body.P("}")

		g.b_body.P()
	}

	return nil
}

// Get type converter for type
func (g *Generator) getTypeConv(fldtype string) TypeConverter {
	tp, scalar := g.GetDepType(fldtype)

	if !scalar {
		for _, tc := range g.TypeConvs {
			for _, src := range tc.GetSources() {
				if src.FilePath == tp.FileDep.FilePath && src.PackageName == tp.FileDep.ProtoFile.PackageName {
					return NewTypeConverterList([]TypeConverter{tc, g.tc_default})
				}
			}
		}
	}

	return g.tc_default
}

// Get dependent type
func (g *Generator) GetDepType(fldtype string) (tp *fdep.DepType, isscalar bool) {
	// check if if scalar
	if _, ok := fproto.ParseScalarType(fldtype); ok {
		isscalar = true
	} else {
		isscalar = false
		var err error

		tp, err = g.filedep.GetType(fldtype)
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}

// Get type for field.
// If pbsource=true, returns the name for the source Golang generated file. Else returns the name
// generated by GoWrap.
func (g *Generator) GetType(fldtype string, pbsource bool) (t string, tp *fdep.DepType, isscalar bool) {
	// check if if scalar
	if st, ok := fproto.ParseScalarType(fldtype); ok {
		t = st.GoType()
		isscalar = true
	} else {
		isscalar = false
		var err error

		tp, err = g.filedep.GetType(fldtype)
		if err != nil {
			log.Fatal(err)
		}

		if !pbsource && tp.FileDep.IsSame(g.filedep) {
			//_ = g.FileDep(tp.FileDep, tp.Alias)
			t = fmt.Sprintf("%s", tp.Name)
		} else {
			falias := g.FileDep(tp.FileDep, tp.Alias, pbsource)
			t = fmt.Sprintf("%s.%s", falias, tp.Name)
		}
	}

	return
}

// Declares a dependency and returns the alias to be used on this file.
func (g *Generator) Dep(imp string, defalias string) string {
	var alias string
	var ok bool
	if alias, ok = g.imports[imp]; ok {
		return alias
	}

	if defalias == "" {
		defalias = path.Base(imp)
	}

	defalias = strings.Replace(defalias, ".", "_", -1)

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

	g.imports[imp] = alias
	return alias
}

// Declares a dependency using a FileDep.
func (g *Generator) FileDep(filedep *fdep.FileDep, defalias string, pbsource bool) string {
	var p string
	if !pbsource && !filedep.IsSame(g.filedep) && filedep.DepType == fdep.DepType_Own {
		p = g.GoWrapPackage(filedep)
	} else {
		p = filedep.GoPackage()
	}
	return g.Dep(p, defalias)
}

// Returns the generated file as a string.
func (g *Generator) String() string {
	p := baseName(g.GoWrapPackage(g.filedep))

	g.b_head.P("package ", p)
	g.b_head.P()
	for i, ia := range g.imports {
		g.b_head.P("import ", ia, ` "`, i, `"`)
	}
	g.b_head.P()

	return g.b_head.String() + g.b_body.String()
}

// Returns the expected output file path and name
func (g *Generator) Filename() string {
	p := g.GoWrapPackage(g.filedep)
	return path.Join(p, strings.TrimSuffix(path.Base(g.filedep.FilePath), path.Ext(g.filedep.FilePath))+".gpb.go")
}

// Returns the wrapped package name.
func (g *Generator) GoWrapPackage(filedep *fdep.FileDep) string {
	if g.PkgSource != nil {
		if p, ok := g.PkgSource.GetPkg(filedep); ok {
			return p
		}
	}

	for _, o := range filedep.ProtoFile.Options {
		if o.Name == "gowrap_package" {
			return o.Value
		}
	}
	for _, o := range filedep.ProtoFile.Options {
		if o.Name == "go_package" {
			return o.Value
		}
	}
	return path.Dir(filedep.FilePath)
}
