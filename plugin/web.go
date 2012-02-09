package plugin

import (
	descriptor "code.google.com/p/goprotobuf/compiler/descriptor"
	"code.google.com/p/goprotobuf/compiler/generator"
)

// GenerateWebStubs is the core of the plugin package.
// It generates an interface based on the ServiceDescriptorProto and an RPC
// client implementation of the interface as well as three helper functions
// to create the Client and Server necessary to utilize the service over
// RPC.
func (p *Plugin) GenerateWebStubs(svc *descriptor.ServiceDescriptorProto) {
	p.webImports = true

	name := generator.CamelCase(*svc.Name)
	_ = name

	/*
	p.P("// ", name, " is an interface satisfied by the generated client and")
	p.P("// which must be implemented by the object wrapped by the server.")
	p.P("type ", name, " interface {")
	p.In()
	for _, m := range svc.Method {
		method := generator.CamelCase(*m.Name)
		iType := p.ObjectNamed(*m.InputType)
		oType := p.ObjectNamed(*m.OutputType)
		p.P(method, "(in *", p.TypeName(iType), ", out *", p.TypeName(oType), ") error")
	}
	p.Out()
	p.P("}")
	p.P()
	p.P("// internal wrapper for type-safe RPC calling")
	p.P("type rpc", name, "Client struct {")
	p.In()
	p.P("*rpc.Client")
	p.Out()
	p.P("}")
	for _, m := range svc.Method {
		method := generator.CamelCase(*m.Name)
		iType := p.ObjectNamed(*m.InputType)
		oType := p.ObjectNamed(*m.OutputType)
		p.P("func (this rpc", name, "Client) ", method, "(in *", p.TypeName(iType), ", out *", p.TypeName(oType), ") error {")
		p.In()
		p.P(`return this.Call("`, name, ".", method, `", in, out)`)
		p.Out()
		p.P("}")
	}
	p.P()
	p.P("// New", name, "Client returns an *rpc.Client wrapper for calling the methods of")
	p.P("// ", name, " remotely.")
	p.P("func New", name, "Client(conn net.Conn) ", name, " {")
	p.In()
	p.P("return rpc", name, "Client{rpc.NewClientWithCodec(plugin.NewClientCodec(conn))}")
	p.Out()
	p.P("}")
	p.P()
	p.P("// Serve", name, " serves the given ", name, " backend implementation on conn.")
	p.P("func Serve", name, "(conn net.Conn, backend ", name, ") error {")
	p.In()
	p.P("srv := rpc.NewServer()")
	p.P(`if err := srv.RegisterName("`, name, `", backend); err != nil {`)
	p.In()
	p.P("return err")
	p.Out()
	p.P("}")
	p.P("srv.ServeCodec(plugin.NewServerCodec(conn))")
	p.P("return nil")
	p.Out()
	p.P("}")
	p.P()
	p.P("// Dial", name, " returns a ", name, " for calling the ", name, " servince at addr (TCP).")
	p.P("func Dial", name, "(addr string) (", name, ", error) {")
	p.In()
	p.P(`conn, err := net.Dial("tcp", addr)`)
	p.P("if err != nil {")
	p.In()
	p.P("return nil, err")
	p.Out()
	p.P("}")
	p.P("return New", name, "Client(conn), nil")
	p.Out()
	p.P("}")
	p.P()
	p.P("// ListenAndServe", name, " serves the given ", name, " backend implementation")
	p.P("// on all connections accepted as a result of listening on addr (TCP).")
	p.P("func ListenAndServe", name, "(addr string, backend ", name, ") error {")
	p.In()
	p.P(`clients, err := net.Listen("tcp", addr)`)
	p.P("if err != nil {")
	p.In()
	p.P("return err")
	p.Out()
	p.P("}")
	p.P("srv := rpc.NewServer()")
	p.P(`if err := srv.RegisterName("`, name, `", backend); err != nil {`)
	p.In()
	p.P("return err")
	p.Out()
	p.P("}")
	p.P("for {")
	p.In()
	p.P("conn, err := clients.Accept()")
	p.P("if err != nil {")
	p.In()
	p.P("return err")
	p.Out()
	p.P("}")
	p.P("go srv.ServeCodec(plugin.NewServerCodec(conn))")
	p.Out()
	p.P("}")
	p.P(`panic("unreachable")`)
	p.Out()
	p.P("}")
	*/
}
