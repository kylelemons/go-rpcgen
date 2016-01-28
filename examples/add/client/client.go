// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package main

import (
	"crypto/rsa"
	"crypto/tls"
	"flag"
	"fmt"
	"log"

	"github.com/bradhe/go-rpcgen/examples/add/addservice"
)

var (
	certDir = flag.String("certdir", "certs", "The directory to load the X509 certificates from")
)

func openTLSClient(ipPort string) (*tls.Conn, error) {

	// Note this loads standard x509 certificates, test keys can be
	// generated with makecert.sh

	log.Printf("Loading certificates from directory: %s\n", *certDir)
	cert, err := tls.LoadX509KeyPair(*certDir+"/server.pem", *certDir+"/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	// InsecureSkipVerify required for unsigned certs with Go1 and later.
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", ipPort, &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	// This shows the public key of the server, we will accept any key, but
	// we could terminate the connection based on the public key if desired.
	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		log.Printf("Client: Server public key is:\n%x\n", v.PublicKey.(*rsa.PublicKey).N)

	}
	// Lets verify behind the doubt that both ends of the connection
	// have completed the handshake and negotiated a SSL connection
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	// All TLS handling has completed, now to pass the connection off to
	// go-rpcgen/protobuf/AddService
	return conn, err
}

func main() {

	var x int32
	var y int32
	flag.Parse()
	conn, err := openTLSClient("127.0.0.1:8000")
	sum := addservice.NewAddServiceClient(conn)
	if err != nil {
		log.Fatalf("dial: %s", err)
	}

	for {
		fmt.Printf("Enter 2 numbers> ")
		fmt.Scan(&x)
		fmt.Scan(&y)
		fmt.Printf("Sending %d and %d\n", x, y)
		in := &addservice.AddMessage{X: &x, Y: &y}
		out := &addservice.SumMessage{}
		if err := sum.Add(in, out); err != nil {
			log.Fatalf("Add failed with: %s", err)
		}
		if out.Z == nil {
			log.Fatalf("Sum failed with no message returned")
		}
		fmt.Printf("Received %d\n\n", *out.Z)
	}
}
