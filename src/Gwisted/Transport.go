package Gwisted

import (
    "fmt"
    "net"
    "sync"
)

type ITransport interface {
    Write(data []byte) error
    WriteSequence(seq [][]byte) error
    LoseConnection()
    GetPeer() net.Addr
    GetHost() net.Addr
}

type Transport struct {
    conn *net.TCPConn
    mu   *sync.Mutex
    protocol IProtocol
}

func NewTransport(conn *net.TCPConn, protocol IProtocol) *Transport {
    return &Transport {
        conn: conn,
        mu: &sync.Mutex{},
        protocol: protocol,
    }
}

func (self *Transport) Write(data []byte) error {
    _, err := self.conn.Write(data)
    return err
}

func (self *Transport) LoseConnection() {
    fmt.Println("lose connection")
    self.conn.Close()
    self.conn = nil
    // FIXME
    self.protocol.ConnectionLost("")
}

func (self *Transport) GetPeer() net.Addr {
    return self.conn.RemoteAddr()
}

func (self *Transport) GetHost() net.Addr {
    return self.conn.LocalAddr()
}

func (self *Transport) WriteSequence(seq [][]byte) error {
    self.mu.Lock()
    defer self.mu.Unlock()

    for _, data := range seq {
        err := self.Write(data)
        if (err != nil) {
            return err;
        }
    }
    return nil;
}
