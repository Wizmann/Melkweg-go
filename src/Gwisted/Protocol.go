package Gwisted

import (
    logging "Logging"
)

type IDataReceivedHandler interface {
    DataReceived(data []byte)
}

type IConnectionMadeHandler interface {
    ConnectionMade(factory IProtocolFactory)

}

type IConnectionLostHandler interface {
    ConnectionLost(reason string)
}

type IProtocol interface {
    IDataReceivedHandler
    IConnectionMadeHandler
    IConnectionLostHandler

    MakeConnection(t ITransport)
    GetTransport() ITransport
    IsConnected() bool
    Start()
}

type Protocol struct {
    Transport ITransport
    Factory   IProtocolFactory
    connected int

    DataReceivedHandler IDataReceivedHandler
    ConnectionMadeHandler IConnectionMadeHandler
    ConnectionLostHandler IConnectionLostHandler
}

func NewProtocol() *Protocol {
    return &Protocol {
        Transport: nil,
        connected: 0,

        DataReceivedHandler: nil,
        ConnectionMadeHandler: nil,
        ConnectionLostHandler: nil,
    }
}

func (self *Protocol) Start() {
    buf := make([]byte, 65536)
    go func() {
        for {
            n, err := self.Transport.GetConnection().Read(buf)
            if (err == nil) {
                self.DataReceived(buf[:n])
            } else {
                self.ConnectionLost(err.Error())
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
    logging.Debug("DataReceived: %x", data)
    if (self.DataReceivedHandler != nil) {
        self.DataReceivedHandler.DataReceived(data)
    } else {
        // pass
    }
}

func (self *Protocol) ConnectionLost(reason string) {
    self.connected = 0
    if (self.ConnectionLostHandler != nil) {
        self.ConnectionLostHandler.ConnectionLost(reason)
    } else {
        // pass
    }
}

func (self *Protocol) IsConnected() bool {
    return self.connected == 1
}

func (self *Protocol) GetTransport() ITransport {
    return self.Transport
}
