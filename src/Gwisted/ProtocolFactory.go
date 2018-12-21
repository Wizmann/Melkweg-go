package Gwisted

import (
    "net"
)

type ProtocolFactory struct {
    protocolBuilder func(tcp *net.TCPConn) IProtocol
}

func NewProtocolFactory(protocolBuilder func(tcp *net.TCPConn) IProtocol) *ProtocolFactory {
    return &ProtocolFactory {
        protocolBuilder: protocolBuilder,
    }
}

func (self *ProtocolFactory) BuildProtocol(tcp *net.TCPConn) IProtocol {
    return self.protocolBuilder(tcp)
}

