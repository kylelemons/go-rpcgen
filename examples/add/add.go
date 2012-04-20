package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/rsa"
	"flag"
	"log"
	"net"

	"github.com/kylelemons/go-rpcgen/examples/add/addservice"
)

var  (
		certDir = flag.String("certdir","certs","The directory to load the X509 certificates from")
)

// Add is the type which will implement the addservice.AddService interface
// and can be called remotely.  In this case, it does not have any state, but
// it could.
type Add struct{}

// Add is the function that can be called remotely.  Note that this can be
// called concurrently, so if the Add structure did have internal state,
// it should be designed for concurrent access.
func (Add) Add(in *addservice.AddMessage, out *addservice.SumMessage) error {
	out.Z = new(int32)
	*out.Z = *in.X + *in.Y
	log.Printf("server: X=%d Y=%d Z=%d", *in.X, *in.Y, *out.Z)
	return nil
}

func handleClient(conn net.Conn) {
	tlscon, ok := conn.(*tls.Conn)
	if ok {
		log.Print("server: conn: type assert to TLS succeedded")
		err := tlscon.Handshake()
		if err != nil {
			log.Fatalf("server: handshake failed: %s", err)
		} else {
			log.Print("server: conn: Handshake completed")
		}
		state := tlscon.ConnectionState()
		// Note we could reject clients if we don't like their public key.		
		for _, v := range state.PeerCertificates {
			log.Printf("Client: Server public key is:\n%x\n",v.PublicKey.(*rsa.PublicKey).N)
//			log.Printf("Server: client cert chain %s", v.Subject.ToRDNSequence())
		}
		// Now that we have completed SSL/TLS 
		addservice.ServeAddService(tlscon, Add{})
	}
}

func serverTLSListen(service string) {

	// Load x509 certificates for our private/public key, makecert.sh will
	// generate them for you.

	log.Printf("Loading certificates from directory: %s\n",*certDir)
	cert, err := tls.LoadX509KeyPair(*certDir+"/server.pem", *certDir+"/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	// Note if we don't tls.RequireAnyClientCert client side certs are ignored.
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert}
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	// Keep this loop simple/fast as to be able to handle new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		// Fire off go routing to handle rest of connection.
		go handleClient(conn)
	}
}

func main() {
	flag.Parse()
	serverTLSListen("0.0.0.0:8000")
}
