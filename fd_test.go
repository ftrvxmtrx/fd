package fd

import (
	"net"
	"os"
	"sync"
	"testing"
)

var SockFilename = "/tmp/sendfd.sock"

func getFD(t *testing.T, w *sync.WaitGroup) {
	defer w.Done()

	c, err := net.Dial("unix", SockFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	sendfdConn := c.(*net.UnixConn)

	var fs []*os.File
	fs, err = Get(sendfdConn, 1, []string{"a file"})
	if err != nil {
		t.Fatal(err)
	}
	f := fs[0]
	defer f.Close()

	b := make([]byte, 64)
	var n int
	n, err = f.Read(b)
	if err != nil {
		t.Fatal(err)
	}

	if n < 1 {
		t.Fatal("failed to read the data")
	}
}

func TestPutGet(t *testing.T) {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	os.Remove(SockFilename)
	l, err := net.Listen("unix", SockFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	var w sync.WaitGroup
	w.Add(1)
	go getFD(t, &w)

	var a net.Conn
	a, err = l.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()

	listenConn := a.(*net.UnixConn)

	if err = Put(listenConn, f); err != nil {
		t.Fatal(err)
	}

	w.Wait()
}
