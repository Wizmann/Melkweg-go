package Melkweg

import (
    "Gwisted"
    proto "github.com/golang/protobuf/proto"
    "time"
)

type ProtocolState int

const (
    READY   = ProtocolState(1)
    RUNNING = ProtocolState(2)
    DONE    = ProtocolState(3)
    ERROR   = ProtocolState(4)
)

type ILineReceivedOnRunningHandler interface {
    LineReceivedOnRunning(packet *MPacket)
}

type ILineReceivedOnReadyHandler interface {
    LineReceivedOnReady(packet *MPacket)
}

type ProtocolBase struct {
    *Gwisted.Int32StringReceiver

    key          string
    iv           []byte
    cipher       ICipher
    peerCipher   ICipher
    state        ProtocolState
    outgoing     map[int]Gwisted.IProtocol
    timeout      int
    timeoutTimer *time.Timer
    config       *Config

    LineReceivedOnRunningHandler ILineReceivedOnRunningHandler
    LineReceivedOnReadyHandler          ILineReceivedOnReadyHandler
}

func NewProtocolBase() *ProtocolBase {
    config := GetConfigInstance()
    p := &ProtocolBase {
        Int32StringReceiver: Gwisted.NewInt32StringReceiver(),
        key: config.GetKey(),
        config: config,
        timeout: config.GetTimeout(),
        iv: DigestBytes(Nonce(19)),
        state: READY,
    }
    p.setTimeout()
    p.LineReceivedHandler = p
    p.cipher = NewAESCipher(p.iv, p.key)
    p.ConnectionMadeHandler = p
    p.ConnectionLostHandler = p
    return p
}

func (self *ProtocolBase) write(packet *MPacket) error {
    data, err := proto.Marshal(packet)
    if (err != nil) {
        return err
    }

    if (self.state > READY) {
        data = self.cipher.Encrypt(data)
    }

    err = self.SendLine(data)
    if (err != nil) {
        return err;
    }

    return nil;
}

func (self *ProtocolBase) setTimeout() {
    self.timeoutTimer = time.AfterFunc(time.Millisecond * time.Duration(self.timeout), self.timeoutConnection)
}

func (self *ProtocolBase) resetTimeout() {
    self.timeoutTimer.Reset(time.Millisecond * time.Duration(self.timeout))
}

func (self *ProtocolBase) LineReceived(data []byte) {
    packet, err := self.parse(data)
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
    if (self.Transport != nil) {
        self.Transport.LoseConnection()
    } else {
        // clear all outgoings
    }
}

func (self *ProtocolBase) parse(data []byte) (*MPacket, error) {
    packet := &MPacket{}
    if (self.state > READY) {
        data = self.peerCipher.Decrypt(data)
    }
    err := proto.Unmarshal(data, packet)
    return packet, err
}

func (self *ProtocolBase) timeoutConnection() {
    log.Error("connection timeout")
    self.handleError()
}

func (self *ProtocolBase) LineReceivedOnReady(packet *MPacket) {
    if (self.LineReceivedOnReadyHandler != nil) {
        panic("LineReceivedOnReadyHandler is nil")
    }

    self.LineReceivedOnReadyHandler.LineReceivedOnReady(packet)
}

func (self *ProtocolBase) LineReceivedOnRunning(packet *MPacket) {
    if (self.LineReceivedOnRunningHandler != nil) {
        panic("LineReceivedOnRunningHandler is nil")
    }

    self.LineReceivedOnRunningHandler.LineReceivedOnRunning(packet)
}
