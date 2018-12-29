package Gwisted

import (
    "errors"
    "net"
    "sync"
)

type ITransport interface {
    Write(data []byte) error
    WriteSequence(seq [][]byte) error
    LoseConnection()
    GetPeer() *net.TCPAddr
    GetHost() *net.TCPAddr
    GetConnection() *net.TCPConn
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
    if (!self.protocol.IsConnected()) {
        return errors.New("Connection is finished or reset")
    }
    _, err := self.conn.Write(data)
    return err
}

func (self *Transport) LoseConnection() {
    log.Debug("lose connection")
    self.conn.Close()
}

func (self *Transport) GetPeer() *net.TCPAddr {
    return self.conn.RemoteAddr().(*net.TCPAddr)
}

func (self *Transport) GetHost() *net.TCPAddr {
    return self.conn.LocalAddr().(*net.TCPAddr)
}

func (self *Transport) GetConnection() *net.TCPConn {
    return self.conn;
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
