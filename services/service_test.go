package services

import (
	"bytes"
	"strings"
	"testing"

	descriptor "code.google.com/p/goprotobuf/compiler/descriptor"
	"code.google.com/p/goprotobuf/compiler/generator"
	"code.google.com/p/goprotobuf/proto"
)

func TestGenerateService(t *testing.T) {
	cases := []struct{
		Service *descriptor.ServiceDescriptorProto
		Output  string
	}{
		{
			Service: &descriptor.ServiceDescriptorProto{
				Name: proto.String("math"),
				Method: []*descriptor.MethodDescriptorProto{
					{
						Name: proto.String("Sqrt"),
						InputType: proto.String("SqrtInput"),
						OutputType: proto.String("SqrtOutput"),
					},
					{
						Name: proto.String("Add"),
						InputType: proto.String("AddInput"),
						OutputType: proto.String("AddOutput"),
					},
				},
			},
			Output: `
// Math is an interface satisfied by the generated client and
// which must be implemented by the object wrapped by the server.
type Math interface {
	Sqrt(in *SqrtInput, out *SqrtOutput) error
	Add(in *AddInput, out *AddOutput) error
}

// internal wrapper for type-safe RPC calling
type rpcMathClient struct {
	*rpc.Client
}
func (this rpcMathClient) Sqrt(in *SqrtInput, out *SqrtOutput) error {
	return this.Call("Math.Sqrt", in, out)
}
func (this rpcMathClient) Add(in *AddInput, out *AddOutput) error {
	return this.Call("Math.Add", in, out)
}

// NewMathClient returns an *rpc.Client wrapper for calling the methods of
// Math remotely.
func NewMathClient(conn net.Conn) Math {
	return rpcMathClient{rpc.NewClientWithCodec(services.NewClientCodec(conn))}
}

// ServeMath serves the given Math backend implementation on conn.
func ServeMath(conn net.Conn, backend Math) error {
	srv := rpc.NewServer()
	if err := srv.RegisterName("Math", backend); err != nil {
		return err
	}
	srv.ServeCodec(services.NewServerCodec(conn))
	return nil
}

// ListenAndServeMath serves the given Math backend implementation
// on all connections accepted as a result of listening on addr (TCP).
func ListenAndServeMath(addr string, backend Math) error {
	clients, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := rpc.NewServer()
	if err := srv.RegisterName("Math", backend); err != nil {
		return err
	}
	for {
		conn, err := clients.Accept()
		if err != nil {
			return err
		}
		go srv.ServeCodec(services.NewServerCodec(conn))
	}
	panic("unreachable")
}
`,
		},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		p := Plugin{compileGen: fakeCompileGen{&generator.Generator{Buffer: buf}}}
		p.GenerateService(c.Service)
		if got, want := buf.String(), strings.TrimSpace(c.Output)+"\n"; got != want {
			t.Fail()
			t.Logf("GenerateService")
			t.Logf("  Input: %s", c.Service)
			t.Logf("  Got:\n%s", got)
			t.Logf("  Want:\n%s", want)
		}
	}
}

type fakeObject string
func (fakeObject) PackageName() string { return "" }
func (fakeObject) TypeName() []string { return nil }
func (fakeObject) File() *descriptor.FileDescriptorProto { return nil }

type fakeCompileGen struct { *generator.Generator }
func (fakeCompileGen) ObjectNamed(name string) generator.Object { return fakeObject(name) }
func (fakeCompileGen) TypeName(obj generator.Object) string { return string(obj.(fakeObject)) }
