package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go handler(c)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	for {
		w := bufio.NewWriter(conn)
		r := bufio.NewReader(conn)
		rw := bufio.NewReadWriter(r, w)
		var buf [1024]byte
		n, err := rw.Read(buf[:])
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("remote addr:", conn.RemoteAddr(), "value:", string(buf[:n]))
		// rw.Write([]byte("hello" + string(buf[:n])))
	}
}
