package Gwisted

import (
    "bytes"
    "fmt"
    log "github.com/sirupsen/logrus"
    "net"
    "testing"
    "time"
)

type MyServerProtocol struct {
    Protocol
}

func (self *MyServerProtocol) DataReceived(data []byte) {
    log.Debug("server data received: ", data)
    if (bytes.Compare(data, []byte("ping")) == 0) {
        self.transport.Write([]byte("pong"))
    } else if (bytes.Compare(data, []byte("exit")) == 0) {
        self.transport.Write([]byte("exit"))
        self.transport.LoseConnection()
    }
}

type MyClientProtocol struct {
    Protocol
}

func (self *MyClientProtocol) ConnectionMade() {
    log.Debug("client connection made")
    self.transport.Write([]byte("ping"))
}

func (self *MyClientProtocol) DataReceived(data []byte) {
    log.Debug("client data received ", data)
    if (bytes.Compare(data, []byte("pong")) == 0) {
        self.transport.Write([]byte("exit"))
    } else if (bytes.Compare(data, []byte("exit")) == 0) {
        self.transport.LoseConnection()
    }
}

func NewMyServerProtocol(port int) *MyServerProtocol {
    p := &MyServerProtocol{}
    p.DataReceivedHandler = p

    l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if (err != nil) {
        log.Fatal("TCP listen error: ", err)
    }
    log.Debug("TCP listen")
    conn, err := l.Accept()
    if (err != nil) {
        log.Fatal("TCP accept error: ", err)
    }
    log.Debug("TCP accept")

    t := NewTransport(conn.(*net.TCPConn), p)

    p.transport = t
    return p
}

func NewMyClientProtocol(port int) *MyClientProtocol{
    p := &MyClientProtocol{}
    p.ConnectionMadeHandler = p
    p.DataReceivedHandler = p

    conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
    if (err != nil) {
        log.Fatal("TCP dail error: ", err)
    }

    t := NewTransport(conn.(*net.TCPConn), p)
    p.transport = t
    p.ConnectionMade()
    return p
}

func TestPingPongProtocol(t *testing.T) {
    var server *MyServerProtocol
    var client *MyClientProtocol
    go func() {
        server = NewMyServerProtocol(10000)
        server.Start()
    }()

    go func() {
        client = NewMyClientProtocol(10000)
        client.Start()
    }()

    time.Sleep(time.Millisecond * 500)

    for i := 0; i <= 3; i++ {
        if (server.connected == 0 && client.connected == 0) {
            break;
        }
        if (i == 3) {
            t.Error();
        }
        time.Sleep(time.Millisecond * 500)
    }
}
