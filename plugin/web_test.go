// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package plugin

import (
	"bytes"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func TestGenerateWebStubs(t *testing.T) {
	cases := []struct {
		Service *descriptor.ServiceDescriptorProto
		Output  string
	}{
		{
			Service: &descriptor.ServiceDescriptorProto{
				Name: proto.String("math"),
				Method: []*descriptor.MethodDescriptorProto{
					{
						Name:       proto.String("Sqrt"),
						InputType:  proto.String("SqrtInput"),
						OutputType: proto.String("SqrtOutput"),
					},
					{
						Name:       proto.String("Add"),
						InputType:  proto.String("AddInput"),
						OutputType: proto.String("AddOutput"),
					},
				},
			},
			Output: `
// MathWeb is the web-based RPC version of the interface which
// must be implemented by the object wrapped by the webrpc server.
type MathWeb interface {
	Sqrt(r *http.Request, in *SqrtInput, out *SqrtOutput) error
	Add(r *http.Request, in *AddInput, out *AddOutput) error
}

// internal wrapper for type-safe webrpc calling
type rpcMathWebClient struct {
	remote   *url.URL
	protocol webrpc.Protocol
}

func (this rpcMathWebClient) Sqrt(in *SqrtInput, out *SqrtOutput) error {
	return webrpc.Post(this.protocol, this.remote, "/Math/Sqrt", in, out)
}

func (this rpcMathWebClient) Add(in *AddInput, out *AddOutput) error {
	return webrpc.Post(this.protocol, this.remote, "/Math/Add", in, out)
}

// Register a MathWeb implementation with the given webrpc ServeMux.
// If mux is nil, the default webrpc.ServeMux is used.
func RegisterMathWeb(this MathWeb, mux webrpc.ServeMux) error {
	if mux == nil {
		mux = webrpc.DefaultServeMux
	}
	if err := mux.Handle("/Math/Sqrt", func(c *webrpc.Call) error {
		in, out := new(SqrtInput), new(SqrtOutput)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Sqrt(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	if err := mux.Handle("/Math/Add", func(c *webrpc.Call) error {
		in, out := new(AddInput), new(AddOutput)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Add(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	return nil
}

// NewMathWebClient returns a webrpc wrapper for calling the methods of Math
// remotely via the web.  The remote URL is the base URL of the webrpc server.
func NewMathWebClient(pro webrpc.Protocol, remote *url.URL) Math {
	return rpcMathWebClient{remote, pro}
}
`,
		},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		p := Plugin{compileGen: fakeCompileGen{&generator.Generator{Buffer: buf}}}
		p.GenerateWebStubs(c.Service)
		if got, want := buf.String(), strings.TrimSpace(c.Output)+"\n"; got != want {
			t.Fail()
			t.Logf("GenerateRPCStubs")
			t.Logf("  Input: %s", c.Service)
			t.Logf("  Got:\n%s", got)
			t.Logf("  Want:\n%s", want)
		}
	}
}
