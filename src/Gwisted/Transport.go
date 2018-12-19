package Gwisted

import (
    "net"
    "sync"
}

type ITransport interface {
    Write(data []byte) error
    WriteSequence(seq [][]byte) error
    LoseConnection()
    GetPeer() Addr
    GetHost() Addr
}

type Transport struct {
    Conn *TCPConn
    mu   *sync.Mutex{}
}

func (self *Transport) Write(data []byte) error {
    _, err = self.Conn.Write(data)
    return err
}

func (self *Transport) LoseConnection() {
    self.Conn.Close()
}

func (self *Transport) GetPeer() Addr {
    return self.Conn.RemoteAddr()
}

func (self *Transport) GetHost() Addr {
    return self.Conn.LocalAddr()
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
