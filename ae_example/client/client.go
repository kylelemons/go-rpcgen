package main

import (
	"os"
	"log"
	"net/url"
	"github.com/kylelemons/go-rpcgen/ae_example/whoami"
	"github.com/kylelemons/go-rpcgen/webrpc"
)

func main() {
	for _, arg := range os.Args[1:] {
		url, err := url.Parse(arg)
		if err != nil {
			log.Printf("invald url %q: %s", arg, err)
			continue
		}

		svc := whoami.NewWhoamiServiceWebClient(webrpc.JSON, url)
		
		in, out := whoami.Empty{}, whoami.YouAre{}
		if err := svc.Whoami(&in, &out); err != nil {
			log.Printf("whoami(%q): %s", url, err)
			continue
		}
		log.Printf("You are %s", *out.IpAddr)
	}
}
