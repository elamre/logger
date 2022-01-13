package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/elamre/logger/internal"
	"github.com/elamre/logger/pkg/logger"
	"log"
	"net"
)

var mainLogger = logger.NewLogger()

func main() {
	p := make([]byte, 2048)
	byteBuf := new(bytes.Buffer)

	dec := gob.NewDecoder(byteBuf)

	conn, err := net.Dial("udp", "127.0.0.1:12321")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	if _, err := conn.Write([]byte("Connect\n")); err != nil {
		panic(err)
	}
	msg := internal.LogMessage{}
	for {
		if n, err := conn.Read(p); err != nil {
			//if _, err = conn.Read(byteBuf.Bytes()); err != nil {
			panic(err)
			//fmt.Printf("err read: %s", err.Error())
		} else {
			byteBuf.Write(p[:n])
		}

		if err := dec.Decode(&msg); err != nil {
			fmt.Printf("Some error %v\n", err)
			//break
		}
		log.Printf("%+v", msg)
		/**/
	}
	_ = conn.Close()
}
