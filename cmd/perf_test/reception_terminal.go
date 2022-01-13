package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/elamre/logger/internal"
	"net"
	"strings"
)

func startTerminal(port int, count int, expectedLevel string, message string) func() {
	p := make([]byte, 2048)
	byteBuf := new(bytes.Buffer)
	msg := internal.LogMessage{}
	cnt := 0
	ready := make(chan bool)

	dec := gob.NewDecoder(byteBuf)

	conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic(err)
	}
	if _, err := conn.Write([]byte("Connect\n")); err != nil {
		panic(err)
	}

	go func() {
		for cnt < count {
			if n, err := conn.Read(p); err != nil {
				panic(err)
			} else {
				byteBuf.Write(p[:n])
			}

			if err := dec.Decode(&msg); err != nil {
				panic(err)
			}
			if resultCheck {
				if !strings.Contains(msg.StringFormat, message) {
					panic(fmt.Errorf("did not find expected \"%s\" in message \"%s\"", message, msg.StringFormat))
				}
			}
			cnt++
		}
		_ = conn.Close()
		ready <- true
	}()
	return func() {
		<-ready
	}
}
