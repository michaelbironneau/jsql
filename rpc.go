package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

type JSQL int

var port = flag.Int("port", 5123, "the port to listen on")

func (s *JSQL) Select(arg *string, reply Rowset) error {
	var (
		err     error
		selArgs SelectArgs
	)

	if arg == nil {
		return nil //got nil args => return nil reply.
	}

	log.Println("received:", *arg)

	err = json.Unmarshal([]byte(*arg), &selArgs)

	if err != nil {
		return nil //TODO: Don't be silent about malformed args...
	}

	reply, err = selArgs.Select()

	return nil
}

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(*port))
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
