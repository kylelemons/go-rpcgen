// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package webrpc

import (
	"encoding/gob"
	"io"
)

// Gob implements the Go Object implementation of the webrpc.Protocol
// interface.
var Gob Protocol = gobProtocol{}

type gobProtocol struct{}

func (gobProtocol) String() string                            { return "application/x-gob" }
func (gobProtocol) Encode(w io.Writer, obj interface{}) error { return gob.NewEncoder(w).Encode(obj) }
func (gobProtocol) Decode(r io.Reader, obj interface{}) error { return gob.NewDecoder(r).Decode(obj) }

func init() {
	RegisterProtocol(Gob)
}
