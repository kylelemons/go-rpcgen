// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package main

import (
	"flag"
	"testing"

	"github.com/bradhe/go-rpcgen/examples/echo/echoservice"
)

var server = flag.String("server", "localhost:9999", "RPC server address")

func TestEcho(t *testing.T) {
	flag.Parse()

	tests := []string{
		"this is a test",
		"woo, more tests\n",
		"",
	}

	go main()

	echo, err := echoservice.DialEchoService(*server)
	if err != nil {
		t.Fatalf("dial: %s", err)
	}

	for _, test := range tests {
		in := &echoservice.Payload{Message: &test}
		out := &echoservice.Payload{}
		if err := echo.Echo(in, out); err != nil {
			t.Fatalf("echo: %s", err)
		}
		if out.Message == nil {
			t.Fatalf("echo: no message returned")
		}
		if got, want := *out.Message, test; got != want {
			t.Errorf("echo(%q) = %q, want %q", test, got, want)
		}
	}
}
