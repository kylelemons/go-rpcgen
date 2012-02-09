package plugin

import (
	descriptor "code.google.com/p/goprotobuf/compiler/descriptor"
	"code.google.com/p/goprotobuf/compiler/generator"
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
