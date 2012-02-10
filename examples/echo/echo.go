package main

import (
	"flag"
	"log"

	"github.com/kylelemons/go-rpcgen/examples/echo/echoservice"
)

var (
	addr    = flag.String("addr", ":9999", "RPC Server address (transient)")
	message = flag.String("message", "test", "Test echo message")
)

// Echo is the type which will implement the echoservice.EchoService interface
// and can be called remotely.  In this case, it does not have any state, but
// it could.
type Echo struct{}

// Echo is the function that can be called remotely.  Note that this can be
// called concurrently, so if the Echo structure did have internal state,
// it should be designed for concurrent access.
func (Echo) Echo(in *echoservice.Payload, out *echoservice.Payload) error {
	out.Message = in.Message
	return nil
}

func main() {
	flag.Parse()

	// Serve requests on the given address.  We are going to run this in the
	// background, since we're going to connect to our own service within the
	// same daemon.  In general, these will be separate.
	if err := echoservice.ListenAndServeEchoService(*addr, Echo{}); err != nil {
		log.Fatalf("listenandserve: %s", err)
	}
}
