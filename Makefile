SHELL 	 = /bin/bash
GO 			 = /usr/local/go/bin/go

build:
	$(GO) build ./...

setup:
	$(GO) get ./...

install:
	$(GO) install -v ./...
