// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// Package plugin implements a plugin for protoc-gen-go that generates
// RPC stubs for use with the the net/rpc package.
//
// To register the plugin, import this package as follows:
//   import _ "github.com/bradhe/go-rpcgen/plugin"
package plugin

import (
	"os"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/generator"
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

	p.stubs = []string{"rpc", "web"}

	if stubs := os.Getenv("GO_STUBS"); stubs != "" {
		p.stubs = strings.Split(stubs, ",")
	}
}

// Generate generates the RPC stubs for all plugin in the given
// FileDescriptorProto.
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	for _, svc := range file.Service {
		p.GenerateCommonStubs(svc)
		for _, stub := range p.stubs {
			switch stub {
			case "rpc":
				p.GenerateRPCStubs(svc)
			case "web":
				p.GenerateWebStubs(svc)
			default:
				p.Fail("unknown go_stub", stub)
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
		p.P(`import "github.com/bradhe/go-rpcgen/codec"`)
	}
	if p.webImports {
		p.P(`import "net/url"`)
		p.P(`import "net/http"`)
		p.P(`import "github.com/bradhe/go-rpcgen/webrpc"`)
	}
}

func init() {
	generator.RegisterPlugin(new(Plugin))
}
