package main

import (
    . "Melkweg"
    Utils "Melkweg/Utils"
)

type MelkwegClientProtocol struct {
    *ProtocolBase

    heartbeatTimer *time.Timer
}

func NewMelkwegClientProtocol() *MelkwegClientProtocol {
    p := &MelkwegClientProtocol {
        Protocol: NewProtocolBase(),
    }
    p.LineReceivedOnRunningHandler = p
    p.LineReceivedOnReadyHandler = p
    p.ConnectionMadeHandler = p
    p.ConnectionLostHandler = p
    return p
}

func (self *MelkwegClientProtocol) ConnectionMade() {
    log.Debug("connection is made")
    self.write(PacketFactory.CreateSynPacket(self.iv))
}

func (self *MelkwegClientProtocol) ConnectionLost() {
    for port, protocol := range self.outgoing { 
        protocol.transport.LoseConnection()
    }
}

func (self *MelkwegClientProtocol) HandleDataPacket(packet) {
    port := packet.GetPort()

    if protocol, ok := self.outgoing[port]; ok {
        protocol.transport.write(packet.GetData())
    } else {
        self.write(PacketFactory.CreateRstPacket(port))
    }
}

func (self *MelkwegClientProtocol) LineReceivedOnReady(packet *MPacket) {
    if (packet.iv != nil) {
        self.peerCipher = NewAESCipher(self.key, packet.iv)
        log.Infof("get iv from: %s from %s", packet.iv, self.transport.getPeer())
        self.state = RUNNING
        self.heartbeat()
        self.resetTimeout()
    } else {
        self.handleError()
    }
}

func (self *MelkwegClientProtocol) LineReceivedOnRunning(packet *MPacket) {
    if (packet.flags == DATA) {
        self.HandleDataPacket(packet)
    } else if (packet.flags == RST || packet.flags == FIN) {
        // will not happen on client side
        log.Warnf("connection on port %d will be terminated", packet.GetPort())
    } else if (packet.flags == LIV) {
        prevTime := Utils.GetTimestamp()
        log.Warnf("[HEARTBEAT] ping = %d ms", Utils.GetTimestamp() - prevTime)
    } else {
        self.handleError()
    }
}

func (self *MelkwegClientProtocol) heartbeat() {
    packet := PacketFactory.CreateLivPacket()
    packet.ClientTime = Utils.GetTimestamp()
    self.write(packet)
    heartbeatTimer = time.AfterFunc(time.Millisecond * time.Duration(timeout), self.heartbeat)
}

