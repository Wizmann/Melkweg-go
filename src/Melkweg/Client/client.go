package main

import (
    "math/rand"
    "net"
    . "Gwisted"
    . "Melkweg"
    "time"
    logging "Logging"
)

type ClientLocalProxyProtocol struct {
    *Protocol
    clientProtocol *MelkwegClientProtocol
}

func (self *ClientLocalProxyProtocol) DataReceived(data []byte) {
    port := self.Transport.GetPeer().Port
    logging.Verbose("get data on local port: %d", port)
    self.clientProtocol.Write(NewDataPacket(port, data))
}

func (self *ClientLocalProxyProtocol) ConnectionLost(reason error) {
    port := self.Transport.GetPeer().Port
    if (reason.Error() == "ConnectionReset") {
        self.clientProtocol.Write(NewRstPacket(port))
    } else {
        self.clientProtocol.Write(NewFinPacket(port))
    }
    self.Transport.LoseConnection()
    self.clientProtocol.RemovePeer(port)
}

func NewClientLocalProxyProtocol(clientProtocol *MelkwegClientProtocol) *ClientLocalProxyProtocol {
    p := &ClientLocalProxyProtocol{
        Protocol: NewProtocol(),
        clientProtocol: clientProtocol,
    }
    p.DataReceivedHandler = p
    p.ConnectionLostHandler = p
    return p
}

type ClientLocalProxyProtocolFactory struct {
    *ProtocolFactory
    outgoing []*MelkwegClientProtocol
}

func NewClientLocalProxyProtocolFactory() *ClientLocalProxyProtocolFactory {
    config := GetConfigInstance()
    logging.Debug("server addr: %s, server port: %d", config.GetServerAddr(), config.GetServerPort())
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
    logging.Debug("build local proxy protocol")
    n := len(self.outgoing)
    if (n <= 0) {
        panic("no outgoing protocol")
    }

    outgoing := self.outgoing[rand.Intn(len(self.outgoing))]

    if (outgoing == nil) {
        panic("outgoing protocol is nil")
    }

    p := NewClientLocalProxyProtocol(outgoing)
    t := NewTransport(tcp, p)
    outgoing.SetPeer(p, t.GetPeer().Port)
    p.MakeConnection(t)
    return p
}

func main() {
    rand.Seed(time.Now().Unix())

    config := GetConfigInstance()

    factory := NewClientLocalProxyProtocolFactory()
    ReactorInstance.ListenTCP(config.GetClientPort(), factory, 50)

    ReactorInstance.Run()
}
