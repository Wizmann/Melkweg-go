package Gwisted

import (
    "sync"
)

type IDataReceivedHandler interface {
    DataReceived(data []byte)
}

type IConnectionMadeHandler interface {
    ConnectionMade(factory IProtocolFactory)
}

type IConnectionLostHandler interface {
    ConnectionLost(reason error)
}

type IProtocol interface {
    IDataReceivedHandler
    IConnectionMadeHandler
    IConnectionLostHandler

    MakeConnection(t ITransport)
    GetTransport() ITransport
    Start()
}

type Protocol struct {
    Transport ITransport
    Factory   IProtocolFactory
    connected int
    Mu        *sync.Mutex
    pauseCh   chan int
    isPaused  bool

    DataReceivedHandler IDataReceivedHandler
    ConnectionMadeHandler IConnectionMadeHandler
    ConnectionLostHandler IConnectionLostHandler
}

func NewProtocol() *Protocol {
    return &Protocol {
        Transport: nil,
        connected: 0,
        Mu: &sync.Mutex{},
        pauseCh: make(chan int, 1),
        isPaused: false,

        DataReceivedHandler: nil,
        ConnectionMadeHandler: nil,
        ConnectionLostHandler: nil,
    }
}

func (self *Protocol) PauseProducing() {
    self.isPaused = true
}

func (self *Protocol) ResumeProducing() {
    self.isPaused = false
    self.pauseCh <- 1
}

func (self *Protocol) Start() {
    buf := make([]byte, 65536)
    go func() {
        for {
            if (self.isPaused) {
                log.Debug("protocol is paused")
                _ = <-self.pauseCh
            }
            n, err := self.Transport.GetConnection().Read(buf)
            if (err == nil) {
                self.DataReceived(buf[:n])
            } else {
                self.ConnectionLost(err)
                break
            }
        }
    }()
}

func (self *Protocol) MakeConnection(transport ITransport) {
    self.connected = 1
    self.Transport = transport
    self.Start()
}

func (self *Protocol) ConnectionMade(factory IProtocolFactory) {
    if (self.ConnectionMadeHandler != nil) {
        self.ConnectionMadeHandler.ConnectionMade(factory)
    } else {
        self.Factory = factory
        // pass
    }
}

func (self *Protocol) DataReceived(data []byte) {
    if (self.DataReceivedHandler != nil) {
        self.DataReceivedHandler.DataReceived(data)
    } else {
        // pass
    }
}

func (self *Protocol) ConnectionLost(reason error) {
    self.connected = 0
    close(self.pauseCh)
    if (self.ConnectionLostHandler != nil) {
        self.ConnectionLostHandler.ConnectionLost(reason)
    } else {
        // pass
    }
}

func (self *Protocol) GetTransport() ITransport {
    return self.Transport
}
