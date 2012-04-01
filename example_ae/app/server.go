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
