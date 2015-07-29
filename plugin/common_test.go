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

func TestGenerateCommonStubs(t *testing.T) {
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
// Math is an interface satisfied by the generated client and
// which must be implemented by the object wrapped by the server.
type Math interface {
	Sqrt(in *SqrtInput, out *SqrtOutput) error
	Add(in *AddInput, out *AddOutput) error
}
`,
		},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		p := Plugin{compileGen: fakeCompileGen{&generator.Generator{Buffer: buf}}}
		p.GenerateCommonStubs(c.Service)
		if got, want := buf.String(), strings.TrimSpace(c.Output)+"\n"; got != want {
			t.Fail()
			t.Logf("GenerateCommonStubs")
			t.Logf("  Input: %s", c.Service)
			t.Logf("  Got:\n%s", got)
			t.Logf("  Want:\n%s", want)
		}
	}
}
