package Gwisted

import (
    "net"
    "errors"
    logging "Logging"
)

type ITransport interface {
    Write(data []byte) error
    LoseConnection()
    GetPeer() *net.TCPAddr
    GetHost() *net.TCPAddr
    GetConnection() *net.TCPConn
}

type Transport struct {
    conn *net.TCPConn
    queue chan []byte
    ctrlCh chan int
    protocol IProtocol
}

func NewTransport(conn *net.TCPConn, protocol IProtocol) *Transport {
    t := &Transport {
        conn: conn,
        queue: make(chan []byte, 100000000),
        ctrlCh: make(chan int, 1000000),
        protocol: protocol,
    }

    go func() {
        for {
            select {
            case data := <- t.queue:
                t.DoWrite(data)
            case <- t.ctrlCh:
                return
            }
        }
    }()
    return t
}

func (self *Transport) DoWrite(data []byte) {
    if (!self.protocol.IsConnected()) {
        logging.Error("Connection is finished or reset")
        return
    }

    p := 0
    len := len(data)

    for p < len {
        delta, err := self.conn.Write(data[p:])
        if (err != nil) {
            logging.Error(err.Error())
            self.LoseConnection()
            return
        }
        p += delta
    }
}

func (self *Transport) Write(data []byte) error {
    if (!self.protocol.IsConnected()) {
        return errors.New("Connection is finished or reset")
    }
    self.queue <- data
    return nil
}

func (self *Transport) LoseConnection() {
    logging.Debug("lose connection")
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

