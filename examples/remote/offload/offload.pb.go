// Code generated by protoc-gen-go.
// source: examples/remote/offload/offload.proto
// DO NOT EDIT!

/*
Package offload is a generated protocol buffer package.

It is generated from these files:
	examples/remote/offload/offload.proto

It has these top-level messages:
	DataSet
	ResultSet
*/
package offload

import proto "github.com/golang/protobuf/proto"
import math "math"

import "net"
import "net/rpc"
import "github.com/bradhe/go-rpcgen/codec"
import "net/url"
import "net/http"
import "github.com/bradhe/go-rpcgen/webrpc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type DataSet struct {
	Data             *string `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DataSet) Reset()         { *m = DataSet{} }
func (m *DataSet) String() string { return proto.CompactTextString(m) }
func (*DataSet) ProtoMessage()    {}

func (m *DataSet) GetData() string {
	if m != nil && m.Data != nil {
		return *m.Data
	}
	return ""
}

type ResultSet struct {
	Result           *string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ResultSet) Reset()         { *m = ResultSet{} }
func (m *ResultSet) String() string { return proto.CompactTextString(m) }
func (*ResultSet) ProtoMessage()    {}

func (m *ResultSet) GetResult() string {
	if m != nil && m.Result != nil {
		return *m.Result
	}
	return ""
}

func init() {
}

// OffloadService is an interface satisfied by the generated client and
// which must be implemented by the object wrapped by the server.
type OffloadService interface {
	Compute(in *DataSet, out *ResultSet) error
}

// internal wrapper for type-safe RPC calling
type rpcOffloadServiceClient struct {
	*rpc.Client
}

func (this rpcOffloadServiceClient) Compute(in *DataSet, out *ResultSet) error {
	return this.Call("OffloadService.Compute", in, out)
}

// NewOffloadServiceClient returns an *rpc.Client wrapper for calling the methods of
// OffloadService remotely.
func NewOffloadServiceClient(conn net.Conn) OffloadService {
	return rpcOffloadServiceClient{rpc.NewClientWithCodec(codec.NewClientCodec(conn))}
}

// ServeOffloadService serves the given OffloadService backend implementation on conn.
func ServeOffloadService(conn net.Conn, backend OffloadService) error {
	srv := rpc.NewServer()
	if err := srv.RegisterName("OffloadService", backend); err != nil {
		return err
	}
	srv.ServeCodec(codec.NewServerCodec(conn))
	return nil
}

// DialOffloadService returns a OffloadService for calling the OffloadService servince at addr (TCP).
func DialOffloadService(addr string) (OffloadService, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewOffloadServiceClient(conn), nil
}

// ListenAndServeOffloadService serves the given OffloadService backend implementation
// on all connections accepted as a result of listening on addr (TCP).
func ListenAndServeOffloadService(addr string, backend OffloadService) error {
	clients, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := rpc.NewServer()
	if err := srv.RegisterName("OffloadService", backend); err != nil {
		return err
	}
	for {
		conn, err := clients.Accept()
		if err != nil {
			return err
		}
		go srv.ServeCodec(codec.NewServerCodec(conn))
	}
	panic("unreachable")
}

// OffloadServiceWeb is the web-based RPC version of the interface which
// must be implemented by the object wrapped by the webrpc server.
type OffloadServiceWeb interface {
	Compute(r *http.Request, in *DataSet, out *ResultSet) error
}

// internal wrapper for type-safe webrpc calling
type rpcOffloadServiceWebClient struct {
	remote   *url.URL
	protocol webrpc.Protocol
}

func (this rpcOffloadServiceWebClient) Compute(in *DataSet, out *ResultSet) error {
	return webrpc.Post(this.protocol, this.remote, "/OffloadService/Compute", in, out)
}

// Register a OffloadServiceWeb implementation with the given webrpc ServeMux.
// If mux is nil, the default webrpc.ServeMux is used.
func RegisterOffloadServiceWeb(this OffloadServiceWeb, mux webrpc.ServeMux) error {
	if mux == nil {
		mux = webrpc.DefaultServeMux
	}
	if err := mux.Handle("/OffloadService/Compute", func(c *webrpc.Call) error {
		in, out := new(DataSet), new(ResultSet)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Compute(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	return nil
}

// NewOffloadServiceWebClient returns a webrpc wrapper for calling the methods of OffloadService
// remotely via the web.  The remote URL is the base URL of the webrpc server.
func NewOffloadServiceWebClient(pro webrpc.Protocol, remote *url.URL) OffloadService {
	return rpcOffloadServiceWebClient{remote, pro}
}
