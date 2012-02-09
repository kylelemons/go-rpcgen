// Package plugin implements a plugin for protoc-gen-go that generates
// RPC stubs for use with the the net/rpc package.
//
// To register the plugin, import this package as follows:
//   import _ "github.com/kylelemons/go-rpcgen/codec"
package plugin

import (
	"strconv"

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
}

// Name returns the name of the plugin.
func (p *Plugin) Name() string { return "go-rpcgen" }

// Init stores the given generator in the Plugin for use in the
// Generate* class of functions.
func (p *Plugin) Init(g *generator.Generator) {
	p.compileGen = g
}

// Generate generates the RPC stubs for all plugin in the given
// FileDescriptorProto.
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	rpcStubs, webStubs := true, false

	if options := file.FileDescriptorProto.Options; options != nil {
		for _, option := range options.UninterpretedOption {
			if len(option.Name) != 1 || option.Name[0].NamePart == nil {
				continue
			}
			var err error
			name := *option.Name[0].NamePart
			switch name {
			case "go_rpc_stubs":
				rpcStubs, err = strconv.ParseBool(string(option.StringValue))
			case "go_web_stubs":
				webStubs, err = strconv.ParseBool(string(option.StringValue))
			}
			if err != nil {
				p.Error(err, "Could not parse value of " + name)
			}
		}
	}

	if rpcStubs {
		for _, svc := range file.Service {
			p.GenerateRPCStubs(svc)
		}
	}

	if webStubs {
		for _, svc := range file.Service {
			p.GenerateWebStubs(svc)
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
}

func init() {
	generator.RegisterPlugin(new(Plugin))
}
