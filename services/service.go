package services

import (
	"bufio"
	"encoding/binary"
	"io"
	"fmt"
	"net"
	"net/rpc"

	descriptor "code.google.com/p/goprotobuf/compiler/descriptor"
	"code.google.com/p/goprotobuf/compiler/generator"
	"code.google.com/p/goprotobuf/proto"

	"github.com/kylelemons/go-rpcgen/services/wire"
)

// TODO: Use io.ReadWriteCloser instead of net.Conn?

func (p *Plugin) GenerateService(svc *descriptor.ServiceDescriptorProto) {
	p.imports = true

	name := generator.CamelCase(*svc.Name)

	p.P("// ", name, " is an interface satisfied by the generated client and")
	p.P("// which must be implemented by the object wrapped by the server.")
	p.P("type ", name, " interface {")
	p.In()
	for _, m := range svc.Method {
		method := generator.CamelCase(*m.Name)
		iType := p.ObjectNamed(*m.InputType)
		oType := p.ObjectNamed(*m.OutputType)
		p.P(method, "(in *", p.TypeName(iType), ", out *", p.TypeName(oType), ") error")
	}
	p.Out()
	p.P("}")
	p.P()
	p.P("// internal wrapper for type-safe RPC calling")
	p.P("type rpc", name, "Client struct {")
	p.In()
	p.P("*rpc.Client")
	p.Out()
	p.P("}")
	for _, m := range svc.Method {
		method := generator.CamelCase(*m.Name)
		iType := p.ObjectNamed(*m.InputType)
		oType := p.ObjectNamed(*m.OutputType)
		p.P("func (this rpc", name, "Client) ", method, "(in *", p.TypeName(iType), ", out *", p.TypeName(oType), ") error {")
		p.In()
		p.P(`return this.Call("`, name, ".", method, `", in, out)`)
		p.Out()
		p.P("}")
	}
	p.P()
	p.P("// New", name, "Client returns an *rpc.Client wrapper for calling the methods of")
	p.P("// ", name, " remotely.")
	p.P("func New", name, "Client(conn net.Conn) ", name, " {")
	p.In()
	p.P("return rpc", name, "Client{rpc.NewClientWithCodec(services.NewClientCodec(conn))}")
	p.Out()
	p.P("}")
	p.P()
	p.P("// Serve", name, " serves the given ", name, " backend implementation on conn.")
	p.P("func Serve", name, "(conn net.Conn, backend ", name, ") error {")
	p.In()
	p.P("srv := rpc.NewServer()")
	p.P(`if err := srv.RegisterName("`, name, `", backend); err != nil {`)
	p.In()
	p.P("return err")
	p.Out()
	p.P("}")
	p.P("srv.ServeCodec(services.NewServerCodec(conn))")
	p.P(`panic("unreachable")`)
	p.Out()
	p.P("}")
}

type ServerCodec struct {
	r *bufio.Reader
	w io.WriteCloser
}

func NewServerCodec(conn net.Conn) *ServerCodec {
	return &ServerCodec{bufio.NewReader(conn), conn}
}

func (s *ServerCodec) ReadRequestHeader(req *rpc.Request) error {
	size, err := binary.ReadUvarint(s.r)
	if err != nil {
		return err
	}
	// TODO max size?
	message := make([]byte, size)
	if _, err := io.ReadFull(s.r, message); err != nil {
		return err
	}
	var header wire.Header
	if err := proto.Unmarshal(message, &header); err != nil {
		return err
	}
	if header.Method == nil {
		return fmt.Errorf("header missing method: %s", header)
	}
	if header.Seq == nil {
		return fmt.Errorf("header missing seq: %s", header)
	}
	req.ServiceMethod = *header.Method
	req.Seq = *header.Seq
	return nil
}

func (s *ServerCodec) ReadRequestBody(pb interface{}) error {
	size, err := binary.ReadUvarint(s.r)
	if err != nil {
		return err
	}
	// TODO max size?
	message := make([]byte, size)
	if _, err := io.ReadFull(s.r, message); err != nil {
		return err
	}
	return proto.Unmarshal(message, pb)
}

func (s *ServerCodec) WriteResponse(resp *rpc.Response, pb interface{}) error {
	var header wire.Header
	var size []byte
	var data []byte
	var err error

	// Allocate enough space for the biggest size
	size = make([]byte, binary.MaxVarintLen64)

	// Write the header
	if resp.Error != "" {
		header.Error = &resp.Error
	}
	header.Method = &resp.ServiceMethod
	header.Seq = &resp.Seq
	if data, err = proto.Marshal(&header); err != nil {
		return err
	}
	size = size[:binary.PutUvarint(size, uint64(len(data)))]
	if _, err = s.w.Write(size); err != nil {
		return err
	}
	if _, err = s.w.Write(data); err != nil {
		return err
	}

	// Write the proto
	size = size[:cap(size)]
	if _, invalid := pb.(rpc.InvalidRequest); invalid {
		data = nil
	} else {
		if data, err = proto.Marshal(pb); err != nil {
			return err
		}
	}
	size = size[:binary.PutUvarint(size, uint64(len(data)))]
	if _, err = s.w.Write(size); err != nil {
		return err
	}
	if _, err = s.w.Write(data); err != nil {
		return err
	}

	// All done
	return nil
}

func (s *ServerCodec) Close() error {
	return s.w.Close()
}

type ClientCodec struct {
	r *bufio.Reader
	w io.WriteCloser
}

func NewClientCodec(conn net.Conn) *ClientCodec {
	return &ClientCodec{bufio.NewReader(conn), conn}
}

func (c *ClientCodec) WriteRequest(req *rpc.Request, pb interface{}) error {
	var header wire.Header
	var size []byte
	var data []byte
	var err error

	// Allocate enough space for the biggest size
	size = make([]byte, binary.MaxVarintLen64)

	// Write the header
	header.Method = &req.ServiceMethod
	header.Seq = &req.Seq
	if data, err = proto.Marshal(&header); err != nil {
		return err
	}
	size = size[:binary.PutUvarint(size, uint64(len(data)))]
	if _, err = c.w.Write(size); err != nil {
		return err
	}
	if _, err = c.w.Write(data); err != nil {
		return err
	}

	// Write the proto
	size = size[:cap(size)]
	if data, err = proto.Marshal(pb); err != nil {
		return err
	}
	size = size[:binary.PutUvarint(size, uint64(len(data)))]
	if _, err = c.w.Write(size); err != nil {
		return err
	}
	if _, err = c.w.Write(data); err != nil {
		return err
	}

	// All done
	return nil
}

func (c *ClientCodec) ReadResponseHeader(resp *rpc.Response) error {
	size, err := binary.ReadUvarint(c.r)
	if err != nil {
		return err
	}
	// TODO max size?
	message := make([]byte, size)
	if _, err := io.ReadFull(c.r, message); err != nil {
		return err
	}
	var header wire.Header
	if err := proto.Unmarshal(message, &header); err != nil {
		return err
	}
	if header.Method == nil {
		return fmt.Errorf("header missing method: %s", header)
	}
	if header.Seq == nil {
		return fmt.Errorf("header missing seq: %s", header)
	}
	resp.ServiceMethod = *header.Method
	resp.Seq = *header.Seq
	if header.Error != nil {
		resp.Error = *header.Error
	}
	return nil
}

func (c *ClientCodec) ReadResponseBody(pb interface{}) error {
	size, err := binary.ReadUvarint(c.r)
	if err != nil {
		return err
	}
	if size == 0 {
		return nil
	}

	// TODO max size?
	message := make([]byte, size)
	if _, err := io.ReadFull(c.r, message); err != nil {
		return err
	}
	return proto.Unmarshal(message, pb)
}

func (c *ClientCodec) Close() error {
	return c.w.Close()
}
