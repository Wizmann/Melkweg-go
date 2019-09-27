package Gwisted

import (
    logging "Logging"
)

type IDataReceivedHandler interface {
    DataReceived(data []byte)
}

type IConnectionMadeHandler interface {
    ConnectionMade()
}

type IConnectionLostHandler interface {
    ConnectionLost(reason string)
}

type IProtocol interface {
    IDataReceivedHandler
    IConnectionMadeHandler
    IConnectionLostHandler

    makeConnection(t ITransport)
    Start()
}

type Protocol struct {
    transport ITransport
    connected int

    DataReceivedHandler IDataReceivedHandler
    ConnectionMadeHandler IConnectionMadeHandler
    ConnectionLostHandler IConnectionLostHandler
}

func NewProtocol() *Protocol {
    return &Protocol {
        transport: nil,
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
            n, err := self.transport.GetConnection().Read(buf)
            if (err == nil) {
                self.DataReceived(buf[:n])
            } else {
                self.ConnectionLost(err.Error())
                break
            }
        }
    }()
}

func (self *Protocol) makeConnection(transport ITransport) {
    self.connected = 1
    self.transport = transport
    self.ConnectionMade()
    self.Start()
}

func (self *Protocol) ConnectionMade() {
    if (self.ConnectionMadeHandler != nil) {
        self.ConnectionMadeHandler.ConnectionMade()
    } else {
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
