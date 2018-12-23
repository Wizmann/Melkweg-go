package Melkweg

import (
    proto "github.com/golang/protobuf/proto"
)

type ProtocolState int

const (
    READY   = ProtocolState(1)
    RUNNING = ProtocolState(2)
    DONE    = ProtocolState(3)
    ERROR   = ProtocolState(4)
)

type ProtocolBase struct {
    key string
    iv []byte
    cipher ICipher
    state ProtocolState
    outgoing map[int]IProtocol
}

func NewProtocolBase() *ProtocolBase {
    config := Config.GetInstance()
    return &ProtocolBase {
        key: config.GetKey(),
        timeout: config.GetTimeout(),
        iv = Nonce(),
        state = READY
    };
}

func (self *ProtocolBase) write(packet *MPacket) error {
    data, err := proto.Marshal(packet)
    if (err != nil) {
        return err
    }

    if (self.state > READY) {
        data = self.aes.Encrypt(data)
    }

    err = self.sendString(data)
    if (err != nil) {
        return err;
    }

    return nil;
}

func (self *ProtocolBase) LineReceived(data []byte) {
    packet, err = self.parse(data)
    if (err != nil) {
        self.handleError()
    }

    switch self.state {
    case READY:
        self.LineReceivedOnReady(packet)
    case RUNNING:
        self.LineReceivedOnRunning(packet)
    default:
        self.handleError()
    }
}

func (self *ProtocolBase) handleError() {
    self.state = ERROR
    if (self.transport != nil) {
        self.tranport.LoseConnection()
    } else {
        // clear all outgoings
    }
}

func (self *ProtocolBase) parse(data []byte) *MPacket, error {
    packet := &MPacket{}
    if (self.state > READY) {
        data = self.peerAes.Decrypt(data)
    }
    err := proto.Unmarshal(data, packet)
    return packet, err
}
