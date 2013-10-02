// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// +build appengine

package server

import (
	"net/http"
	"whoami"

	_ "github.com/kylelemons/go-rpcgen/webrpc"
)

type server struct{}

func (server) Whoami(r *http.Request, _ *whoami.Empty, out *whoami.YouAre) error {
	out.IpAddr = &r.RemoteAddr
	return nil
}

func init() {
	whoami.RegisterWhoamiServiceWeb(server{}, nil)
}
