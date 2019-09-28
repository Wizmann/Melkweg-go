package Gwisted

import (
    "net"
)

type ProtocolFactory struct {
    protocolBuilder func(tcp *net.TCPConn) IProtocol
}

func ProtocolFactoryForProtocol(protocolCtor func() IProtocol) *ProtocolFactory {
    return &ProtocolFactory {
        protocolBuilder: func(tcp *net.TCPConn) IProtocol {
            p := protocolCtor()
            t := NewTransport(tcp, p)
            p.MakeConnection(t)
            return p
        },
    }
}

func NewProtocolFactory(protocolBuilder func(tcp *net.TCPConn) IProtocol) *ProtocolFactory {
    return &ProtocolFactory {
        protocolBuilder: protocolBuilder,
    }
}

func (self *ProtocolFactory) BuildProtocol(tcp *net.TCPConn) IProtocol {
    return self.protocolBuilder(tcp)
}

