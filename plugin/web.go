// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package plugin

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

// GenerateWebStubs generates the webrpc stubs.
// It generates a webrpc client implementation of the interface as well as
// webrpc handlers and a helper function.
func (p *Plugin) GenerateWebStubs(svc *descriptor.ServiceDescriptorProto) {
	p.webImports = true

	name := generator.CamelCase(*svc.Name)

	p.P(`// `, name, `Web is the web-based RPC version of the interface which`)
	p.P(`// must be implemented by the object wrapped by the webrpc server.`)
	p.P(`type `, name, `Web interface {`)
	p.In()
	for _, m := range svc.Method {
		method := generator.CamelCase(*m.Name)
		iType := p.ObjectNamed(*m.InputType)
		oType := p.ObjectNamed(*m.OutputType)
		p.P(method, "(r *http.Request, in *", p.TypeName(iType), ", out *", p.TypeName(oType), ") error")
	}
	p.Out()
	p.P(`}`)
	p.P()
	p.P("// internal wrapper for type-safe webrpc calling")
	p.P("type rpc", name, "WebClient struct {")
	p.In()
	p.P("remote   *url.URL")
	p.P("protocol webrpc.Protocol")
	p.Out()
	p.P("}")
	p.P()
	for _, m := range svc.Method {
		method := generator.CamelCase(*m.Name)
		iType := p.ObjectNamed(*m.InputType)
		oType := p.ObjectNamed(*m.OutputType)
		p.P("func (this rpc", name, "WebClient) ", method, "(in *", p.TypeName(iType), ", out *", p.TypeName(oType), ") error {")
		p.In()
		p.P(`return webrpc.Post(this.protocol, this.remote, "/`, name, "/", method, `", in, out)`)
		p.Out()
		p.P("}")
		p.P()
	}
	p.P(`// Register a `, name, `Web implementation with the given webrpc ServeMux.`)
	p.P(`// If mux is nil, the default webrpc.ServeMux is used.`)
	p.P(`func Register`, name, `Web(this `, name, `Web, mux webrpc.ServeMux) error {`)
	p.In()
	p.P(`if mux == nil {`)
	p.In()
	p.P(`mux = webrpc.DefaultServeMux`)
	p.Out()
	p.P(`}`)
	for _, m := range svc.Method {
		method := generator.CamelCase(*m.Name)
		iType := p.ObjectNamed(*m.InputType)
		oType := p.ObjectNamed(*m.OutputType)
		p.P(`if err := mux.Handle("/`, name, "/", method, `", func(c *webrpc.Call) error {`)
		p.In()
		p.P(`in, out := new(`, p.TypeName(iType), `), new(`, p.TypeName(oType), `)`)
		p.P(`if err := c.ReadRequest(in); err != nil {`)
		p.In()
		p.P(`return err`)
		p.Out()
		p.P(`}`)
		p.P(`if err := this.`, method, `(c.Request, in, out); err != nil {`)
		p.In()
		p.P(`return err`)
		p.Out()
		p.P(`}`)
		p.P(`return c.WriteResponse(out)`)
		p.Out()
		p.P("}); err != nil {")
		p.In()
		p.P("return err")
		p.Out()
		p.P("}")
	}
	p.P("return nil")
	p.Out()
	p.P("}")
	p.P()
	p.P("// New", name, "WebClient returns a webrpc wrapper for calling the methods of ", name)
	p.P("// remotely via the web.  The remote URL is the base URL of the webrpc server.")
	p.P("func New", name, "WebClient(pro webrpc.Protocol, remote *url.URL) ", name, " {")
	p.In()
	p.P("return rpc", name, "WebClient{remote, pro}")
	p.Out()
	p.P("}")
}
