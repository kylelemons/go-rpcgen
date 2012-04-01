package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kylelemons/go-rpcgen/examples/remote/offload"
	"github.com/kylelemons/go-rpcgen/webrpc"
)

var addr = flag.String("addr", ":9999", "RPC server bind address")

type OffloadService struct{}

func (o *OffloadService) Compute(r *http.Request, in *offload.DataSet, out *offload.ResultSet) error {
	if in.Data == nil {
		return nil
	}

	str := *in.Data
	res := make([]byte, len(str))
	last := len(str) - 1
	for i := range str {
		res[last-i] = str[i]
	}
	str = string(res)

	out.Result = &str
	return nil
}

func main() {
	flag.Parse()

	offload.RegisterOffloadServiceWeb(&OffloadService{}, nil)
	if err := webrpc.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("listenandserve: %s", err)
	}
}
