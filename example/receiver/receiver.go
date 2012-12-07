package main

import (
	"flag"
	"github.com/goerlang/fd"
	"io"
	"log"
	"net"
	"os"
)

var (
	socket string
)

func init() {
	flag.StringVar(&socket, "s", "/tmp/sendfd.sock", "socket")
}

func main() {
	flag.Parse()

	if !flag.Parsed() || socket == "" {
		flag.Usage()
		os.Exit(1)
	}

	c, err := net.Dial("unix", socket)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	fdConn := c.(*net.UnixConn)

	var fs []*os.File
	fs, err = fd.Get(fdConn, 1, []string{"a file"})
	if err != nil {
		log.Fatal(err)
	}
	f := fs[0]
	defer f.Close()

	for {
		b := make([]byte, 4096)
		var n int
		n, err = f.Read(b)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s", b[:n])
	}
}
