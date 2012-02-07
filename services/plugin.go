package services

import (
	"code.google.com/p/goprotobuf/compiler/generator"
)

var _ generator.Plugin = &Plugin{}

type compileGen interface {
	// Output
	P(...interface{})
	In()
	Out()

	// Object lookup
	ObjectNamed(string) generator.Object
	TypeName(generator.Object) string
}

// Plugin implements the generator.Plugin interface
type Plugin struct {
	imports    bool
	compileGen
}

func (p *Plugin) Name() string { return "go-rpcgen" }

func (p *Plugin) Init(g *generator.Generator) {
	p.compileGen = g
}

func (p *Plugin) Generate(file *generator.FileDescriptor) {
	for _, svc := range file.Service {
		p.GenerateService(svc)
	}
}

func (p *Plugin) GenerateImports(file *generator.FileDescriptor) {
	if !p.imports {
		return
	}
	p.P("import (")
	p.In()
	p.P(`"net"`)
	p.P(`"net/rpc"`)
	p.P()
	p.P(`"github.com/kylelemons/go-rpcgen/services"`)
	p.Out()
	p.P(")")
}

func init() {
	generator.RegisterPlugin(new(Plugin))
}
