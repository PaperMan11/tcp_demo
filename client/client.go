package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"
)

func NormalFactory() (*MyConn, error) {
	c, err := net.Dial("tcp", ":1234")
	if err != nil {
		return nil, err
	}
	return &MyConn{Conn: c}, nil
}

type Connection interface {
	io.Closer
	io.Writer
	// Write(b []byte) (int, error)
}

type MyConn struct {
	Conn net.Conn
}

var _ Connection = (*MyConn)(nil)

func (c *MyConn) Write(b []byte) (n int, err error) {
	w := bufio.NewWriterSize(c.Conn, 1024)
	w.Write(b)
	w.Flush()
	return
}

func (c *MyConn) Close() error {
	return c.Conn.Close()
}

func main() {
	// c, err := net.Dial("tcp", ":1234")
	// if err != nil {
	// 	panic(err)
	// }
	// defer c.Close()

	// w := bufio.NewWriterSize(c, 1024)
	// w.Write([]byte("hello"))
	// w.Flush()
	// //time.Sleep(time.Second)
	p, err := NewPool(NormalFactory, 5)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	var c *MyConn
	c, err = p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	c.Write([]byte("111"))

	p.Release(c)
	time.Sleep(time.Second)
	c2, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	c2.Write([]byte("222"))
}
