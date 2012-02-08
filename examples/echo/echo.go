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
	go echoservice.ListenAndServeEchoService(*addr, Echo{})

	// Dial the EchoService (which is ourselves in this example)
	e, err := echoservice.DialEchoService(*addr)
	if err != nil {
		log.Fatalf("dial: %s", err)
	}

	var in, out echoservice.Payload
	in.Message = message
	if err := e.Echo(&in, &out); err != nil {
		log.Fatalf("echo: %s", err)
	}
	log.Printf("echo: %s", *out.Message)
}
