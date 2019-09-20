package Gwisted

import (
    "net"
    "sync"
    logging "Logging"
)

type ITransport interface {
    Write(data []byte) error
    WriteSequence(seq [][]byte) error
    LoseConnection()
    GetPeer() net.Addr
    GetHost() net.Addr
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
    _, err := self.conn.Write(data)
    return err
}

func (self *Transport) LoseConnection() {
    logging.Debug("lose connection")
    self.conn.Close()
}

func (self *Transport) GetPeer() net.Addr {
    return self.conn.RemoteAddr()
}

func (self *Transport) GetHost() net.Addr {
    return self.conn.LocalAddr()
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
