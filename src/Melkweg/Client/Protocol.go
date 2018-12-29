package main

import (
    "encoding/hex"
    . "Gwisted"
    . "Melkweg"
    "time"
    Utils "Melkweg/Utils"
)

type MelkwegClientProtocol struct {
    *ProtocolBase

    heartbeatTimeout int
    heartbeatTimer *time.Timer
}

func NewMelkwegClientProtocol() IProtocol {
    p := &MelkwegClientProtocol {
        ProtocolBase: NewProtocolBase(),
        heartbeatTimeout: GetConfigInstance().GetHeartbeatTimeout(),
    }
    p.LineReceivedOnRunningHandler = p
    p.LineReceivedOnReadyHandler = p
    p.ConnectionMadeHandler = p
    p.ConnectionLostHandler = p
    return p
}

func (self *MelkwegClientProtocol) ConnectionMade(factory IProtocolFactory) {
    self.Factory = factory
    log.Infof("send iv: %s", hex.EncodeToString(self.Iv))
    self.Write(NewSynPacket(self.Iv))
}

func (self *MelkwegClientProtocol) ConnectionLost(reason error) {
    log.Errorf("connectionLost because %s", reason)
    self.heartbeatTimer.Stop()
    self.TimeoutTimer.Stop()
    for _, protocol := range self.Outgoing {
        protocol.GetTransport().LoseConnection()
    }

    if (self.Factory != nil) {
        self.Factory.ClientConnectionLost(reason);
    }
}

func (self *MelkwegClientProtocol) HandleDataPacket(packet *MPacket) {
    port := int(packet.GetPort())
    log.Debugf("packet is for port: %d", port)

    if protocol, ok := self.Outgoing[port]; ok {
        protocol.GetTransport().Write(packet.GetData())
    } else {
        self.Write(NewRstPacket(port))
    }
}

func (self *MelkwegClientProtocol) LineReceivedOnReady(packet *MPacket) {
    if (packet.GetIv() != nil) {
        self.PeerCipher = NewAESCipher(packet.GetIv(), self.Key)
        log.Infof("get iv: %s from %s", hex.EncodeToString(packet.GetIv()), self.GetTransport().GetPeer())
        self.State = RUNNING
        self.ResetTimeout()
        self.heartbeat()
    } else {
        self.HandleError()
    }
}

func (self *MelkwegClientProtocol) LineReceivedOnRunning(packet *MPacket) {
    log.Debugf("packet received on running: %s", PacketToString(packet))
    if (packet.GetFlags() == DATA) {
        self.HandleDataPacket(packet)
    } else if (packet.GetFlags() == RST || packet.GetFlags() == FIN) {
        // will not happen on client side
        log.Warningf("connection on port %d will be terminated", packet.GetPort())
    } else if (packet.GetFlags() == LIV) {
        prevTime := Utils.GetTimestamp()
        log.Warningf("[HEARTBEAT] ping = %d ms", Utils.GetTimestamp() - prevTime)
        self.ResetTimeout()
    } else {
        self.HandleError()
        return
    }
}

func (self *MelkwegClientProtocol) heartbeat() {
    packet := NewLivPacket()
    packet.ClientTime = Utils.GetTimestamp()
    self.Write(packet)
    self.heartbeatTimer = time.AfterFunc(time.Millisecond * time.Duration(self.heartbeatTimeout), self.heartbeat)
}

