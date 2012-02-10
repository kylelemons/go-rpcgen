package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"

	"github.com/kylelemons/go-rpcgen/examples/echo/echoservice"
)

var server = flag.String("server", "localhost:9999", "RPC server address")

func main() {
	echo, err := echoservice.DialEchoService(*server)
	if err != nil {
		log.Fatalf("dial: %s", err)
	}

	lines := bufio.NewReader(os.Stdin)
	for {
		os.Stdout.WriteString("> ")
		line, err := lines.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("read: %s", err)
		}

		in := &echoservice.Payload{Message:&line}
		out := &echoservice.Payload{}
		if err := echo.Echo(in, out); err != nil {
			log.Fatalf("echo: %s", err)
		}
		if out.Message == nil {
			log.Fatalf("echo: no message returned")
		}
		os.Stdout.WriteString("< ")
		os.Stdout.WriteString(*out.Message)
	}
}
