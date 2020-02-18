package main

import (
    "encoding/hex"
    . "Gwisted"
    . "Melkweg"
    "time"
    "sync"
    Utils "Melkweg/Utils"
    logging "Logging"
)

type Map = sync.Map

type MelkwegClientProtocol struct {
    *MelkwegProtocolBase

    config *ProxyConfig
    heartbeatTimeout int
    heartbeatTimer *time.Timer
}

func NewMelkwegClientProtocol(config *ProxyConfig) IProtocol {
    p := &MelkwegClientProtocol {
        config: config,
        MelkwegProtocolBase: NewMelkwegProtocolBase(config),
        heartbeatTimeout: config.GetHeartbeatTimeout(),
    }
    p.LineReceivedOnRunningHandler = p
    p.LineReceivedOnReadyHandler = p
    p.ConnectionMadeHandler = p
    p.ConnectionLostHandler = p
    return p
}

func (self *MelkwegClientProtocol) ConnectionMade(factory IProtocolFactory) {
    self.Factory = factory
    logging.Info("send iv: %s", hex.EncodeToString(self.Iv))
    self.Write(NewSynPacket(self.Iv))
    self.Cipher = CipherFactory(self.CipherType, self.Iv, self.Key)
}

func (self *MelkwegClientProtocol) ConnectionLost(reason string) {
    logging.Error("connectionLost because %s", reason)
    if (self.heartbeatTimer != nil) {
        self.heartbeatTimer.Stop()
    }

    if (self.TimeoutTimer != nil) {
        self.TimeoutTimer.Stop()
    }

    for _, protocol := range self.Outgoing {
        transport := protocol.GetTransport()
        if (transport != nil) {
            transport.LoseConnection()
        }
    }

    if (self.Factory != nil) {
        self.Factory.ClientConnectionLost(reason);
    }
}

func (self *MelkwegClientProtocol) HandleDataPacket(packet *MPacket) {
    port := int(packet.GetPort())
    logging.Verbose("packet is for port: %d", port)

    if protocol, ok := self.Outgoing[port]; ok {
        protocol.GetTransport().Write(packet.GetData())
    } else {
        self.Write(NewRstPacket(port))
    }
}

func (self *MelkwegClientProtocol) LineReceivedOnReady(packet *MPacket) {
    if (packet.GetIv() != nil) {
        self.PeerCipher = CipherFactory(self.CipherType, packet.GetIv(), self.Key)
        logging.Info("get iv: %s from %s", hex.EncodeToString(packet.GetIv()), self.GetTransport().GetPeer())
        self.State = RUNNING
        self.ResetTimeout()
        self.heartbeat()
    } else {
        self.HandleError()
    }
}

func (self *MelkwegClientProtocol) LineReceivedOnRunning(packet *MPacket) {
    logging.Verbose("packet received on running: %s", PacketToString(packet))
    if (packet.GetFlags() == DATA) {
        self.HandleDataPacket(packet)
    } else if (packet.GetFlags() == RST || packet.GetFlags() == FIN) {
        port := int(packet.GetPort())
        if protocol, ok := self.Outgoing[port]; ok {
            protocol.GetTransport().LoseConnection()
            logging.Info("connection on port %d will be terminated", packet.GetPort())
            self.RemovePeer(port)
        }
    } else if (packet.GetFlags() == LIV) {
        prevTime := packet.GetClientTime()
        logging.Info("[%s - HEARTBEAT] ping = %d ms", self.config.GetName(), Utils.GetTimestamp() - prevTime)
        logging.Debug("[HEARTBEAT] get ping = %d, server = %d, now = %d",
                packet.GetClientTime(),
                packet.GetServerTime(),
                Utils.GetTimestamp());
        self.ResetTimeout()
    } else {
        self.HandleError()
        return
    }
}

func (self *MelkwegClientProtocol) heartbeat() {
    packet := NewLivPacket()
    packet.ClientTime = Utils.GetTimestamp()
    before_t := Utils.GetTimestamp()
    self.Write(packet)
    after_t := Utils.GetTimestamp()
    logging.Debug("[HEARTBEAT] send ping = %d, time spent: %d", packet.ClientTime, after_t - before_t);
    self.heartbeatTimer = time.AfterFunc(time.Millisecond * time.Duration(self.heartbeatTimeout), self.heartbeat)
}
