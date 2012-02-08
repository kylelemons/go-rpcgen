// Package services implements a plugin for protoc-gen-go that generates
// RPC stubs for use with the the net/rpc package.
//
// To register the plugin, import this package as follows:
//   import _ "github.com/kylelemons/go-rpcgen/services"
package services

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

	// Object lookup
	ObjectNamed(string) generator.Object
	TypeName(generator.Object) string
}

// Plugin implements the generator.Plugin interface.
type Plugin struct {
	imports bool
	compileGen
}

// Name returns the name of the plugin.
func (p *Plugin) Name() string { return "go-rpcgen" }

// Init stores the given generator in the Plugin for use in the
// Generate* class of functions.
func (p *Plugin) Init(g *generator.Generator) {
	p.compileGen = g
}

// Generate generates the RPC stubs for all services in the given
// FileDescriptorProto.
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	for _, svc := range file.Service {
		p.GenerateService(svc)
	}
}

// GenerateImports adds the required imports to the output file if the Generate
// function generated any RPC stubs.
func (p *Plugin) GenerateImports(file *generator.FileDescriptor) {
	if !p.imports {
		return
	}
	p.P(`import "net"`)
	p.P(`import "net/rpc"`)
	p.P(`import "github.com/kylelemons/go-rpcgen/services"`)
}

func init() {
	generator.RegisterPlugin(new(Plugin))
}
