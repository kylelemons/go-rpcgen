// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package plugin

import (
	descriptor "code.google.com/p/goprotobuf/protoc-gen-go/descriptor"
	"code.google.com/p/goprotobuf/protoc-gen-go/generator"
)

type fakeObject string

func (fakeObject) PackageName() string                   { return "" }
func (fakeObject) TypeName() []string                    { return nil }
func (fakeObject) File() *descriptor.FileDescriptorProto { return nil }

type fakeCompileGen struct{ *generator.Generator }

func (fakeCompileGen) ObjectNamed(name string) generator.Object { return fakeObject(name) }
func (fakeCompileGen) TypeName(obj generator.Object) string     { return string(obj.(fakeObject)) }
