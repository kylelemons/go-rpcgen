// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/bradhe/go-rpcgen/examples/remote/offload"
	"github.com/bradhe/go-rpcgen/webrpc"
)

var base = flag.String("base", "http://localhost:9999/", "RPC server base URL")

func main() {
	url, err := url.Parse(*base)
	if err != nil {
		log.Fatalf("url: %s", err)
	}

	oops := "oops, something bad happened"
	do := func(pro webrpc.Protocol, s string) string {
		in := &offload.DataSet{Data: &s}
		out := &offload.ResultSet{Result: &oops}

		off := offload.NewOffloadServiceWebClient(pro, url)
		if err := off.Compute(in, out); err != nil {
			log.Fatalf("compute(%q) - %s", s, err)
		}
		return *out.Result
	}

	log.Printf(do(webrpc.JSON, "I'm a reversed JSON request"))
	log.Printf(do(webrpc.ProtoBuf, "I am sent as a protobuf"))
}
