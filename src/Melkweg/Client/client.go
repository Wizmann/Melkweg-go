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

func (self *ClientLocalProxyProtocol) ConnectionLost(reason string) {
    port := self.Transport.GetPeer().Port
    if (reason == "ConnectionReset") {
        self.clientProtocol.Write(NewRstPacket(port))
    } else {
        self.clientProtocol.Write(NewFinPacket(port))
    }
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
    Config *ProxyConfig
    outgoing []*MelkwegClientProtocol
}

func NewClientLocalProxyProtocolFactory(config ProxyConfig) *ClientLocalProxyProtocolFactory {
    logging.Debug("server addr: %s, server port: %d", config.GetServerAddr(), config.GetServerPort())
    n := config.GetClientOutgoingConnectionNum()

    f := &ClientLocalProxyProtocolFactory {}
    f.Config = &config
    f.ProtocolFactory = NewProtocolFactory(f.BuildProtocol)

    f.outgoing = make([]*MelkwegClientProtocol, n)

    for i := 0; i < n; i++ {
        idx := i
        protocolBuilder := func() IProtocol {
            p := NewMelkwegClientProtocol(f.Config)
            f.outgoing[idx] = p.(*MelkwegClientProtocol)
            return p
        }
        factory := NewReconnectingClientProtocolFactoryForProtocol(
            protocolBuilder, 500, 500, 5000, 2.0, -1)

        Reactor.ConnectTCP(
            config.GetServerAddr(),
            config.GetServerPort(),
            factory,
            -1)
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

    for _, proxyConfig := range config.ProxyConfigs {
        factory := NewClientLocalProxyProtocolFactory(proxyConfig)
        logging.Info("start TCP for Client on port %d", proxyConfig.GetClientPort())
        Reactor.ListenTCP(proxyConfig.GetClientPort(), factory, 50)
    }

    Reactor.Start()
}
