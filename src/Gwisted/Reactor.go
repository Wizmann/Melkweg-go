package Gwisted

import (
    "fmt"
    "net"
)

type Reactor struct {
    ctrlCh chan int
}

func newReactor() *Reactor {
    return &Reactor {
        ctrlCh: make(chan int),
    }
}

func (self *Reactor) Run() {
    for {
        select {
        case <- self.ctrlCh:
            break
        }
    }
}

func (self *Reactor) Stop() {
    self.ctrlCh <- -1
}

func (self *Reactor) ListenTCP(port int, factory *ProtocolFactory, backlog int) {
    go func() {
        l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
        if (err != nil) {
            log.Errorf("listen TCP on port %d error: %s!", port, err.Error())
            return
        }
        for {
            conn, err := l.Accept()
            if (err != nil) {
                log.Errorf("accept TCP on port %d error!", port)
                continue
            }
            _ = factory.BuildProtocol(conn.(*net.TCPConn))
        }
    }()
}

func (self *Reactor) ConnectTCP(host string, port int, factory *ProtocolFactory, timeout int) (IProtocol, error) {
    c := NewClientConnector(host, port, timeout)
    factory.SetConnector(c)
    conn, err := c.Connect()
    if (err != nil) {
        factory.ClientConnectionLost(err)
        return nil, err
    }
    return factory.BuildProtocol(conn), nil
}

var reactor = newReactor()
