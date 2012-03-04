Introduction
============

There are some problems using protobuf with appengine; notably, its use of unsafe.
Here are some directions about how to use go-rpcgen on AppEngine.

Deploying to AppEngine
======================

When deploying to AppEngine, you need two things:A

1. The `go-rpcgen/webrpc` package
1. Sanitized .pb.go files

I have provided scripts in the root of the repository to help out with these.
The `ae_install.sh` script should be run from your appengine project's root directory
(the one containing app.yaml) and will copy the webrpc package into the proper place.
This will also automatically delete the local copy of proto.go,
which removes the protobuf support and the dependency upon it.
The `ae_protoc.go` script (which relies on protoc-gen-go being in your PATH)
will compile (for Go only) all of the .proto files specified on the command-line
with support for web services only (to avoid the dependency on the codec package)
and will sanitize the generated file to remove references to goprotobuf.
A side-effect of the sanitization is that there will no longer be a .String()
method on the generated objects; you may add one manually if you wish,
but I would recommend doing it in a parallel .go file so that regeneration won't kill it.

Testing This Example
--------------------
- Update app.yaml with your application name and the latest Go SDK Version
- Execute the following to run locally:

        ../ae_install.sh
        ../aex_protoc.sh whoami/*.proto
        (cd github.com/kylelemons/go-rpcgen/; mkdir -p ae_example; ln -s ../../../../whoami ae_example/)
        dev_appserver.py .
        go run client/client.go http://localhost:6060/

- Run the following to test remotely:

        appcfg.py update .
        go run client/client.go http://path-to-your-app.appspot.com/
