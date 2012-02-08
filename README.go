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
//   go get -v -u --fix github.com/kylelemons/go-rpcgen/protoc-gen-go
//
// The --fix option is (as of 2012-02-07) required to fix the imports of the
// goprotobuf source to use the new googlecode layout.  You will need to run
// "go fix" on the output of this package before the .pb.go will compile under
// the go tool.
//
// Usage
//
// Usage of the package is pretty straightforward.  Once you have installed the
// protoc-gen-go plugin, you can compile protobufs with the following command
// (where file.proto is the protocol buffer file(s) in question):
//   protoc --go_out=. file.proto && go fix
//
// This will generate a file named like file.pb.go which contains, in addition
// to the usual Go bindings for the messages, an interface for each service
// containing the methods for that service and functions for creating and using
// them with the RPC package.  As mentioned above, the "go fix" is necessary to
// fix the imports in the generated .pb.go file.
//
// See the examples/ subdirectory for some simple examples demonstrating basic
// usage.
package documentation
