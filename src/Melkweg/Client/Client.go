package main

import (
    "math/rand"
    "net"
    . "Gwisted"
    . "Melkweg"
    "time"
)

type ClientLocalProxyProtocol struct {
    *Protocol
    clientProtocol *MelkwegClientProtocol
}

func (self *ClientLocalProxyProtocol) DataReceived(data []byte) {
    port := self.Transport.GetPeer().Port
    log.Debugf("get data on local port: %d", port)
    self.clientProtocol.Write(NewDataPacket(port, data))
}

func (self *ClientLocalProxyProtocol) ConnectionLost(reason string) {
    port := self.Transport.GetPeer().Port
    if (reason == "ConnectionReset") {
        self.clientProtocol.Write(NewRstPacket(port))
    } else {
        self.clientProtocol.Write(NewFinPacket(port))
    }
    self.clientProtocol.RemovePeer(port)
}

func NewClientLocalProxyProtocol(clientProtocol *MelkwegClientProtocol, conn *net.TCPConn) *ClientLocalProxyProtocol {
    p := &ClientLocalProxyProtocol{
        Protocol: NewProtocol(),
        clientProtocol: clientProtocol,
    }
    t := NewTransport(conn, p)
    p.DataReceivedHandler = p
    p.ConnectionLostHandler = p
    p.MakeConnection(t)
    return p
}

type ClientLocalProxyProtocolFactory struct {
    *ProtocolFactory
    outgoing []*MelkwegClientProtocol
}

func NewClientLocalProxyProtocolFactory() *ClientLocalProxyProtocolFactory {
    config := GetConfigInstance()
    log.Debugf("server addr: %s, server port: %d", config.GetServerAddr(), config.GetServerPort())
    n := config.GetClientOutgoingConnectionNum()

    f := &ClientLocalProxyProtocolFactory {}
    f.ProtocolFactory = NewProtocolFactory(f.BuildProtocol)

    f.outgoing = make([]*MelkwegClientProtocol, n)
    for i := 0; i < n; i++ {
        factory := NewReconnectingClientProtocolFactoryForProtocol(
            NewMelkwegClientProtocol, 500, 500, 5000, 2.0, -1)

        p, _ := ReactorInstance.ConnectTCP(
            config.GetServerAddr(),
            config.GetServerPort(),
            factory,
            -1)
        f.outgoing[i] = p.(*MelkwegClientProtocol)
    }
    return f
}

func (self *ClientLocalProxyProtocolFactory) BuildProtocol(tcp *net.TCPConn) IProtocol {
    log.Debug("build local proxy protocol")
    n := len(self.outgoing)
    if (n <= 0) {
        panic("no outgoing protocol")
    }

    outgoing := self.outgoing[rand.Intn(len(self.outgoing))]

    if (outgoing == nil) {
        panic("outgoing protocol is nil")
    }

    p := NewClientLocalProxyProtocol(outgoing, tcp)
    outgoing.SetPeer(p, p.GetTransport().GetPeer().Port)
    return p
}

func main() {
    rand.Seed(time.Now().Unix())

    factory := NewClientLocalProxyProtocolFactory()
    ReactorInstance.ListenTCP(20011, factory, 50)

    ReactorInstance.Run()
}
