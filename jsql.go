package main

import (
	"crypto/tls"
	"errors"
	"flag"
	jsql "github.com/michaelbironneau/jsql/lib"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

type JSQL int

var (
	port       = flag.Int("port", 5123, "the port to listen on")
	secret     = flag.String("password", "", "password to require from clients (optional)")
	certFile   = flag.String("cert", "", "server certificate for TLS (optional)")
	keyFile    = flag.String("key", "", "server private key for TLS (optional)")
	skipVerify = flag.Bool("skip-verify", false, "skip certificate verification (default: false)")
)

func (s *JSQL) Select(arg *jsql.SelectArgs, reply *jsql.Rowset) error {
	var (
		err error
	)

	if arg == nil {
		return nil //got nil args => return nil reply.
	}

	if secret != nil && *secret != arg.Auth {
		return errors.New("incorrect password")
	}

	*reply, err = arg.Select()

	return err
}

func main() {

	flag.Parse()
	var (
		tlsConfig *tls.Config
		l         net.Listener
		err       error
	)

	if len(*certFile) > 0 && len(*keyFile) > 0 {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: *skipVerify,
		}

		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(*certFile, *keyFile)

		if err != nil {
			log.Fatalf("Error loading certificate/keyfile: %s", err.Error())
		}

		l, err = tls.Listen("tcp", "127.0.0.1:"+strconv.Itoa(*port), tlsConfig)

	} else {
		l, err = net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(*port))
	}

	defer l.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Print("listening:", l.Addr())

	jsql := new(JSQL)
	rpc.Register(jsql)

	for {
		log.Print("waiting for connections...")
		c, err := l.Accept()

		if err != nil {
			log.Printf("accept error: %s", c)
			continue
		}

		log.Printf("connection started: %v", c.RemoteAddr())
		go jsonrpc.ServeConn(c)
	}
}
