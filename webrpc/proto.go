// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// +build !appengine

package webrpc

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

// ProtoBuf implements the Google Protocol Buffer implementation of the
// webrpc.Protocol interface.  This is currently not allowed on AppEngine.
var ProtoBuf Protocol = pbProtocol{}

type pbProtocol struct{}

func (pbProtocol) String() string { return "application/protobuf" }
func (pbProtocol) Decode(r io.Reader, obj interface{}) error {
	pb, ok := obj.(proto.Message)
	if !ok {
		return fmt.Errorf("%T does not implement proto.Message", obj)
	}

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("webrpc: decode: %s", err)
	}
	return proto.Unmarshal(body, pb)
}
func (pbProtocol) Encode(w io.Writer, obj interface{}) error {
	pb, ok := obj.(proto.Message)
	if !ok {
		return fmt.Errorf("%T does not implement proto.Message", obj)
	}

	body, err := proto.Marshal(pb)
	if err != nil {
		return fmt.Errorf("webrpc: encode: %s", err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("webrpc: encode: %s", err)
	}
	return nil
}

func init() {
	RegisterProtocol(ProtoBuf)
}
