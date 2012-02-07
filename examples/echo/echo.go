package main

import (
	"flag"
	"log"
	"net"

	"github.com/kylelemons/go-rpcgen/examples/echo/echoservice"
)

var (
	addr = flag.String("addr", ":9999", "RPC Server address (transient)")
	message = flag.String("message", "test", "Test echo message")
)

type Echo struct {}
func (Echo) Echo(in *echoservice.Payload, out *echoservice.Payload) error {
	out.Message = in.Message
	return nil
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen: %s", err)
	}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept: %s")
		}
		echoservice.ServeEchoService(conn, Echo{})
	}()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("dial: %s", err)
	}
	e := echoservice.NewEchoServiceClient(conn)

	var in, out echoservice.Payload
	in.Message = message
	if err := e.Echo(&in, &out); err != nil {
		log.Fatalf("echo: %s", err)
	}
	log.Printf("echo: %s", *out.Message)
}
