package Gwisted

import (
    "net"
)

type IDataReceivedHandler interface {
    DataReceived(data []byte)
}

type IConnectionMadeHandler interface {
    ConnectionMade()
}

type Protocol interface {
    transport *Transport
    connected int

    DataReceivedHandler IDataReceivedHandler
    ConnectionMadeHandler IConnectionMadeHandler
}

func NewProtocol() *Protocol {
    return &Protocol {
        transport: nil,
        connected: 0

        DataReceivedHandler: nil,
        ConnectionMadeHandler: nil,
    }
}

func (self *Protocol) makeConnection(transport *Transport) {
    self.connected = 1
    self.transport = transport
    self.ConnectionMade()
}

func (self *Protocol) ConnectionMade() {
    if (self.ConnectionMadeHandler != nil) {
        self.ConnectionMadeHandler.ConnectionMade()
    } else {
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

