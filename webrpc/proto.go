// +build !appengine

package webrpc

import (
	"fmt"
	"io"
	"io/ioutil"

	"code.google.com/p/goprotobuf/proto"
)

// ProtoBuf implements the Google Protocol Buffer implementation of the
// webrpc.Protocol interface.  This is currently not allowed on AppEngine.
var ProtoBuf Protocol = pbProtocol{}

type pbProtocol struct{}

func (pbProtocol) String() string { return "application/protobuf" }
func (pbProtocol) Decode(r io.Reader, pb interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("webrpc: decode: %s", err)
	}
	return proto.Unmarshal(body, pb)
}
func (pbProtocol) Encode(w io.Writer, pb interface{}) error {
	body, err := proto.Marshal(pb)
	if err != nil {
		return fmt.Errorf("webrpc: encode: %s", err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("webrpc: encode: %s", err)
	}
	return nil
}

func init() {
	RegisterProtocol(ProtoBuf)
}
