package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/kylelemons/go-rpcgen/examples/remote/offload"
	"github.com/kylelemons/go-rpcgen/webrpc"
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
