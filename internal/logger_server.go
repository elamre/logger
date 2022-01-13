package internal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
)

type LogMessage struct {
	TimeNano     int64
	LoggerName   string
	FileName     string
	StringFormat string
}

type LoggerServer struct {
	port              int
	enabled           bool
	logger            LoggerInterface
	conn              *net.UDPConn
	connectedTerminal *net.UDPAddr

	enc *gob.Encoder
	dec *gob.Decoder

	bufMessage LogMessage

	outputBuffer *bytes.Buffer
	iunputBuffer bytes.Buffer
}

func NewLoggerServer(port int) LoggerServer {
	return LoggerServer{
		port: port,
	}
}

func (server *LoggerServer) SetLogger(logger LoggerInterface) {
	server.logger = logger
}

func (server *LoggerServer) Disable() {
	server.enabled = false
	_ = server.conn.Close()
	server.connectedTerminal = nil
}

func (server LoggerServer) IsEnabled() bool {
	return server.enabled
}

func (server LoggerServer) IsConnected() bool {
	return server.enabled && (server.connectedTerminal != nil)
}

func (server *LoggerServer) Run() error {
	var err error
	if server.port < 1000 {
		return fmt.Errorf("use port higher than 1000")
	}
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", server.port))
	if err != nil {
		return err
	}
	server.conn, err = net.ListenUDP("udp4", s)
	if err != nil {
		return err
	}
	server.enabled = true
	server.bufMessage = LogMessage{}
	buffer := make([]byte, 1024)
	name, _ := os.Hostname()
	server.logger.LogInfof("udp log server started at %s:%d", name, server.port)

	go func() {
		defer server.conn.Close()
		for server.enabled {
			n, addr, err := server.conn.ReadFromUDP(buffer)
			if err != nil {
				break
			}
			if string(buffer[:n-1]) == "Connect" {
				server.logger.LogDebug("Client connected")
				// Prevent possible rare case of race condition
				// https://stackoverflow.com/questions/21968266/handling-read-write-udp-connection-in-go
				newAddr := new(net.UDPAddr)
				*newAddr = *addr
				newAddr.IP = make(net.IP, len(addr.IP))
				copy(newAddr.IP, addr.IP)
				server.connectedTerminal = newAddr
				server.outputBuffer = new(bytes.Buffer)
				server.enc = gob.NewEncoder(server.outputBuffer)
				server.dec = gob.NewDecoder(server.conn)
			}
		}
	}()
	return nil
}

func (server *LoggerServer) SendMessageGob(now time.Time, loggerName string, fileName string, format string) error {
	// server.bufMessage.TimeNano = now.UnixNano()
	// server.bufMessage.LoggerName = loggerName
	// server.bufMessage.FileName = fileName
	// server.bufMessage.StringFormat = format
	//server.outputBuffer.Reset()

	if err := server.enc.Encode(&LogMessage{
		TimeNano:     now.UnixNano(),
		LoggerName:   loggerName,
		FileName:     fileName,
		StringFormat: format,
	}); err != nil {
		return err
	}
	output := server.outputBuffer.Bytes()

	// log.Printf("out (%d): % 02x", len(output), output)
	if _, err := server.conn.WriteToUDP(output, server.connectedTerminal); err != nil {
		return err
	}
	server.outputBuffer.Reset()
	return nil
}

func (server *LoggerServer) SendMessage(message string) error {
	if !server.enabled || server.connectedTerminal == nil {
		return nil
	}
	_, err := server.conn.WriteToUDP([]byte(message), server.connectedTerminal)
	return err
}
