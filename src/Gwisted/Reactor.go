package Gwisted

import (
    "fmt"
    "net"
    logging "Logging"
)

type ReactorImpl struct {
    ctrlCh chan int
}

func (self *ReactorImpl) Start() {
    for {
        select {
        case <- self.ctrlCh:
            return
        }
    }
}

func (self *ReactorImpl) Stop() {
    self.ctrlCh <- -1
}

func (self *ReactorImpl) ListenTCP(port int, factory IProtocolFactory, backlog int) {
    go func() {
        l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
        if (err != nil) {
            logging.Fatal("listen TCP on port %d error: %s!", port, err.Error())
            return
        }
        for {
            conn, err := l.Accept()
            if (err != nil) {
                logging.Fatal("accept TCP on port %d error!", port)
                continue
            }
            _ = factory.BuildProtocol(conn.(*net.TCPConn))
        }
    }()
}

func (self *ReactorImpl) ConnectTCP(host string, port int, factory IProtocolFactory, timeout int) (IProtocol, error) {
    c := NewClientConnector(host, port, timeout)
    factory.SetConnector(c)
    conn, err := c.Connect()
    if (err != nil) {
        factory.ClientConnectionLost(err.Error())
        return nil, err
    }
    return factory.BuildProtocol(conn), nil
}

func newReactorImpl() *ReactorImpl {
    return &ReactorImpl {
        ctrlCh: make(chan int),
    }
}

var Reactor = newReactorImpl()
