package main

import (
	"flag"
	"github.com/goerlang/fd"
	"log"
	"net"
	"os"
)

var (
	filename string
	socket   string
)

func init() {
	flag.StringVar(&filename, "f", "", "filename")
	flag.StringVar(&socket, "s", "/tmp/sendfd.sock", "socket")
}

func main() {
	flag.Parse()

	if !flag.Parsed() || filename == "" || socket == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	l, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	var a net.Conn
	a, err = l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	listenConn := a.(*net.UnixConn)
	if err = fd.Put(listenConn, f); err != nil {
		log.Fatal(err)
	}
}
