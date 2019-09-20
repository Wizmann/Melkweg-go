package Gwisted

import (
    "fmt"
    "net"
    logging "Logging"
)

type Reactor struct {
    ctrlCh chan int
}

func (self *Reactor) Start() {
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

func (self *Reactor) ConnectTCP(host string, port int, factory *ProtocolFactory, timeout int) (IProtocol, error) {
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
    if (err != nil) {
        logging.Fatal("dial TCP error on ", host, ":", port)
        return nil, err
    }
    return factory.BuildProtocol(conn.(*net.TCPConn)), nil
}
