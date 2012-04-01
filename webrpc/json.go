package webrpc

import (
	"encoding/json"
	"io"
)

// JSON implements the Javascript Object Notation implementation of the
// webrpc.Protocol interface.
var JSON Protocol = jsonProtocol{}

type jsonProtocol struct{}

func (jsonProtocol) String() string                            { return "application/json" }
func (jsonProtocol) Encode(w io.Writer, obj interface{}) error { return json.NewEncoder(w).Encode(obj) }
func (jsonProtocol) Decode(r io.Reader, obj interface{}) error { return json.NewDecoder(r).Decode(obj) }

func init() {
	RegisterProtocol(JSON)
}
