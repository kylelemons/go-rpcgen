package webrpc

import (
	"encoding/gob"
	"io"
)

// Gob implements the Go Object implementation of the webrpc.Protocol
// interface.
var Gob Protocol = gobProtocol{}

type gobProtocol struct{}

func (gobProtocol) String() string                            { return "application/x-gob" }
func (gobProtocol) Encode(w io.Writer, obj interface{}) error { return gob.NewEncoder(w).Encode(obj) }
func (gobProtocol) Decode(r io.Reader, obj interface{}) error { return gob.NewDecoder(r).Decode(obj) }

func init() {
	RegisterProtocol(Gob)
}
