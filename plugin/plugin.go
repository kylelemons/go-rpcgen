// Package plugin implements a plugin for protoc-gen-go that generates
// RPC stubs for use with the the net/rpc package.
//
// To register the plugin, import this package as follows:
//   import _ "github.com/kylelemons/go-rpcgen/codec"
package plugin

import (
	"code.google.com/p/goprotobuf/compiler/generator"
)

// Fail to compile if Plugin doesn't implement the generator.Plugin interface
var _ generator.Plugin = &Plugin{}

type compileGen interface {
	// Output
	P(...interface{})
	In()
	Out()

	// Errors
	Fail(...string)
	Error(error, ...string)

	// Object lookup
	ObjectNamed(string) generator.Object
	TypeName(generator.Object) string
}

// Plugin implements the generator.Plugin interface.
type Plugin struct {
	rpcImports bool
	webImports bool
	compileGen

	stubs []string
}

// Name returns the name of the plugin.
func (p *Plugin) Name() string { return "go-rpcgen" }

// Init stores the given generator in the Plugin for use in the
// Generate* class of functions.
func (p *Plugin) Init(g *generator.Generator) {
	p.compileGen = g

	// TODO: Figure out some way to derive these
	// - Command-line?
	// - .proto directives?
	p.stubs = []string{"rpc", "web"}
}

// Generate generates the RPC stubs for all plugin in the given
// FileDescriptorProto.
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	for _, stub := range p.stubs {
		for _, svc := range file.Service {
			switch stub {
			case "rpc": p.GenerateRPCStubs(svc)
			case "web": p.GenerateWebStubs(svc)
			default:    p.Fail("unknown go_stub", stub)
			}
		}
	}
}

// GenerateImports adds the required imports to the output file if the Generate
// function generated any RPC stubs.
func (p *Plugin) GenerateImports(file *generator.FileDescriptor) {
	if p.rpcImports {
		p.P(`import "net"`)
		p.P(`import "net/rpc"`)
		p.P(`import "github.com/kylelemons/go-rpcgen/codec"`)
	}
	if p.webImports {
		p.P(`import "net/url"`)
		p.P(`import "net/http"`)
		p.P(`import "github.com/kylelemons/go-rpcgen/webrpc"`)
	}
}

func init() {
	generator.RegisterPlugin(new(Plugin))
}
