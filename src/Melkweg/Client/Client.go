package main

import (
    "math/rand"
    "net"
    "Gwisted"
    "Melkweg"
)

type ClientLocalProxyProtocol struct {
    Protocol

    outgoing IProtocol
}

func NewClientLocalProxyProtocol(outgoing IProtocol, conn *net.TCPConn) IProtocol {
    p := &ClientLocalProxyProtocol{
        outgoing: outgoing
    }
    return p
}

func (self *ClientLocalProxyProtocol) DataReceived(data []byte) {
    port := self.Transport.GetHost().port
    self.outgoing.Write(Melkweg.PacketFactory.CreateDataPacket(port, data))
}

type ClientLocalProxyProtocolFactory struct {
    outgoing []MelkwegClientProtocol
}

func NewClientLocalProxyProtocolFactory() *ClientLocalProxyProtocolFactory {
    config := Melkweg.GetConfigInstance()
    n := config.GetClientOutgoingConnectionNum
    self.outgoing = make([]MelkwegClientProtocol, n)
    for i := 0; i < n; i++ {
        factory := Melkweg.NewReconnectingClientProtocolFactoryForProtocol(
            NewMelkwegClientProtocol, 500, 500, 5000, 2.0, -1)

        outgoing[i] = Gwisted.Reactor.ConnectTCP(
            config.GetServerAddr(),
            config.GetServerPort(),
            factory)
    }
}

func (self *ClientLocalProxyProtocolFactory) BuildProtocol(tcp *net.TCPConn) IProtocol {
    n := len(self.outgoing)
    if (n <= 0) {
        log.Error("no outgoing protocol")
        panic()
    }

    outgoing := self.outgoing[rand.Intn(len(self.outgoing))]

    if (outgoing == nil) {
        panic("outgoing protocol is nil")
    }

    return NewClientLocalProxyProtocol(outgoing, tcp)
}

func main() {
    rand.Seed(time.Now().Unix())

    factory := NewClientLocalProxyProtocolFactory()
    reactor.ListenTCP("20011", factory.BuildProtocol)
    reactor.Run()
}
