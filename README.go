// The go-rpcgen project is an attempt to create an easy-to-use, open source
// protobuf service binding for the standard Go RPC package.  It provides a
// protoc-gen-go (based on the standard "main" from goprotobuf and leveraging
// its libraries) which has a plugin added to also output RPC stub code.
//
// Prerequisites
//
// You will need the protobuf compiler for your operating system of choice.
// You can retrieve this from http://code.google.com/p/protobuf/downloads/list
// if you do not have it already.  As this package builds a plugin for the
// protoc from that package, you will need to have your $GOPATH/bin in your
// path when you run protoc.
//
// Installation
//
// To install, run the following command:
//   go get -v -u github.com/kylelemons/go-rpcgen/protoc-gen-go
//
// Usage
//
// Usage of the package is pretty straightforward.  Once you have installed the
// protoc-gen-go plugin, you can compile protobufs with the following command
// (where file.proto is the protocol buffer file(s) in question):
//   protoc --go_out=. file.proto
//
// This will generate a file named like file.pb.go which contains, in addition
// to the usual Go bindings for the messages, an interface for each service
// containing the methods for that service and functions for creating and using
// them with the RPC package and a webrpc package.
//
// Generated Code for RPC
//
// Given the following basic .proto definition:
//  
//   package echoservice;
//   message payload {
//       required string message = 1;
//   }
//   service echo_service {
//       rpc echo (payload) returns (payload);
//   }
//
// The protoc-gen-go plugin will generate a service definition similar to below:
//   
//   // EchoService is an interface satisfied by the generated client and
//   // which must be implemented by the object wrapped by the server.
//   type EchoService interface {
//       Echo(in *Payload, out *Payload) error
//   }
//
//   // DialEchoService returns a EchoService for calling the EchoService servince at addr (TCP).
//   func DialEchoService(addr string) (EchoService, error) {
//   
//   // NewEchoServiceClient returns an *rpc.Client wrapper for calling the methods of
//   // EchoService remotely.
//   func NewEchoServiceClient(conn net.Conn) EchoService
//   
//   // ListenAndServeEchoService serves the given EchoService backend implementation
//   // on all connections accepted as a result of listening on addr (TCP).
//   func ListenAndServeEchoService(addr string, backend EchoService) error
//
//   // ServeEchoService serves the given EchoService backend implementation on conn.
//   func ServeEchoService(conn net.Conn, backend EchoService) error
//
// Any type which implements EchoService can thus be registered via ServeEchoService
// or ListenAndServeEchoService to be called remotely via NewEchoServiceClient.
//
// Generated Code for WebRPC
//
// In addition to the above, the following are also generated to facilitate
// serving RPCs over the web (e.g. AppEngine):
//   
//   // EchoServiceWeb is the web-based RPC version of the interface which
//   // must be implemented by the object wrapped by the webrpc server.
//   type EchoServiceWeb interface {
//     Echo(r *http.Request, in *Payload, out *Payload) error
//   }
//   
//   // NewEchoServiceWebClient returns a webrpc wrapper for calling the methods of EchoService
//   // remotely via the web.  The remote URL is the base URL of the webrpc server.
//   func NewEchoServiceWebClient(remote *url.URL, pro webrpc.Protocol) EchoService
//
//   // Register a EchoServiceWeb implementation with the given webrpc ServeMux.
//   // If mux is nil, the default webrpc.ServeMux is used.
//   func RegisterEchoServiceWeb(this EchoServiceWeb, mux webrpc.ServeMux) error
//
// Any type which implements EchoServiceWeb (notice that the handlers also
// receive the *http.Request) can be registered.  The RegisterEchoServiceWeb
// function registers the given backend implementation to be called from the
// web via the webrpc package.
//
// Examples
//
// See the examples/ subdirectory for some complete examples demonstrating
// basic usage.
package documentation
