package Gwisted

import (
    "net"
)

type IProtocolFactory interface {
    BuildProtocol(tcp *net.TCPConn) IProtocol
    SetConnector(c IConnector)
    ClientConnectionLost(reason string)
    ClientConnectionFailed(reason string)
}

type IClientConnectionLostHandler interface {
    ClientConnectionLost(reason string)
}

type IClientConnectionFailedHandler interface {
    ClientConnectionFailed(reason string)
}

type ProtocolFactory struct {
    protocolBuilder func(tcp *net.TCPConn) IProtocol
    ClientConnectionLostHandler IClientConnectionLostHandler
    ClientConnectionFailedHandler IClientConnectionFailedHandler

    connector IConnector
    host string
    port int
}

func ProtocolFactoryForProtocol(protocolCtor func() IProtocol) *ProtocolFactory {
    f := &ProtocolFactory {}
    f.protocolBuilder = func(tcp *net.TCPConn) IProtocol {
        p := protocolCtor()
        t := NewTransport(tcp, p)
        p.MakeConnection(t)
        p.ConnectionMade(f)
        return p
    }
    return f
}

func NewProtocolFactory(protocolBuilder func(tcp *net.TCPConn) IProtocol) *ProtocolFactory {
    p := &ProtocolFactory {
        protocolBuilder: protocolBuilder,
    }
    p.ClientConnectionLostHandler = p
    p.ClientConnectionFailedHandler = p
    return p
}

func (self *ProtocolFactory) BuildProtocol(tcp *net.TCPConn) IProtocol {
    return self.protocolBuilder(tcp)
}

func (self *ProtocolFactory) SetConnector(c IConnector) {
    self.connector = c
}

func (self *ProtocolFactory) ClientConnectionLost(reason string) {
    if (self.ClientConnectionLostHandler != nil) {
        self.ClientConnectionLostHandler.ClientConnectionLost(reason)
    }
}

func (self *ProtocolFactory) ClientConnectionFailed(reason string) {
    if (self.ClientConnectionFailedHandler != nil) {
        self.ClientConnectionFailedHandler.ClientConnectionFailed(reason)
    }
}
