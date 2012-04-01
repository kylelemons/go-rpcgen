Introduction
============

There are some problems using protobuf with appengine; notably, its use of unsafe.
Here are some directions about how to use go-rpcgen on AppEngine.

Deploying to AppEngine
======================

When deploying to AppEngine, you need two things:

1. The `go-rpcgen/webrpc` package
1. Sanitized .pb.go files

I have provided scripts in the root of the repository to help out with these.

Getting `go-rpcgen/webrpc` in your app
--------------------------------------

The `ae_install.sh` script should be run from your appengine project's root
directory (the one containing `app.yaml`) and will copy the webrpc package into
the proper place.  The webrpc package contains proto support, but the go1
version of appengine will ignore this file.

Getting sanitized protobufs in your app
---------------------------------------

The `ae_protoc.go` script (which relies on `protoc-gen-go` being in your
`PATH`) will compile (for Go only) all of the `.proto` files specified on the
command-line.

The script has been changed since its previous version to generate both a full
`.pb.go` for use in normal applications and a `.ae.go` for use in appengine.
They both have compilation guards which cause them to be executed in the proper
context.

Both proto files are compiled with support for web services only (to avoid the
dependency on the codec package) and the script will sanitize the duplicate
`.ae.go` file to remove references to goprotobuf.  A side-effect of the
sanitization is that there will no longer be a `.String()` method on the
generated objects under appengine; you may add one manually if you wish, but I
would recommend doing it in an adjacent `.go` file (with the proper `+build`
guard for appengine) so that regeneration won't kill it.

Testing This Example
====================

To test this out, you will probably need `go1` installed locally so that you
can compile non-appengine binaries.  The appengine `go` wrapper script in
particular doesn't like arguments to `go run` though it may be possible to use.

Here are the basic steps:

- Update app.yaml with your application name and the latest Go SDK Version
- Execute the following (in the `example_ae` directory) to run locally:

        ../ae_install.sh
        ../ae_protoc.sh whoami/*.proto
        dev_appserver.py .

- Then, in another shell (since `dev_appserver` blocks):

        go run client/client.go http://localhost:6060/

- Run the following to test remotely:

        appcfg.py update .
        go run client/client.go http://path-to-your-app.appspot.com/
