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

// GenerateCommonStubs is the core of the plugin package.
// It generates an interface based on the ServiceDescriptorProto that is used
// by the other two halves of the plugin.
func (p *Plugin) GenerateCommonStubs(svc *descriptor.ServiceDescriptorProto) {
	name := generator.CamelCase(*svc.Name)

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
}
