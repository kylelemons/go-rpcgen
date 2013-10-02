// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package webrpc

import (
	"encoding/json"
	"io"
)

// JSON implements the Javascript Object Notation implementation of the
// webrpc.Protocol interface.
var JSON Protocol = jsonProtocol{}

type jsonProtocol struct{}

func (jsonProtocol) String() string                            { return "application/json" }
func (jsonProtocol) Encode(w io.Writer, obj interface{}) error { return json.NewEncoder(w).Encode(obj) }
func (jsonProtocol) Decode(r io.Reader, obj interface{}) error { return json.NewDecoder(r).Decode(obj) }

func init() {
	RegisterProtocol(JSON)
}
