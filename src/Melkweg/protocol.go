package Melkweg

import (
    "encoding/hex"
    "Gwisted"
    proto "github.com/golang/protobuf/proto"
    "sync"
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

    Key          string
    Iv           []byte
    Cipher       ICipher
    PeerCipher   ICipher
    State        ProtocolState
    Outgoing     map[int]Gwisted.IProtocol
    Timeout      int
    timeoutTimer *time.Timer
    config       *Config
    mu           *sync.Mutex

    LineReceivedOnRunningHandler  ILineReceivedOnRunningHandler
    LineReceivedOnReadyHandler    ILineReceivedOnReadyHandler
}

func NewProtocolBase() *ProtocolBase {
    config := GetConfigInstance()
    p := &ProtocolBase {
        Int32StringReceiver: Gwisted.NewInt32StringReceiver(),
        Key: config.GetKey(),
        config: config,
        Timeout: config.GetTimeout(),
        Iv: DigestBytes(Nonce(19)),
        State: READY,
        Outgoing: map[int]Gwisted.IProtocol{},
        mu: &sync.Mutex{},
    }
    p.setTimeout()
    p.LineReceivedHandler = p
    p.Cipher = NewAESCipher(p.Iv, p.Key)
    p.ConnectionMadeHandler = p
    p.ConnectionLostHandler = p
    return p
}

func (self *ProtocolBase) Write(packet *MPacket) error {
    mu.Lock()
    defer mu.Unlock()

    data, err := proto.Marshal(packet)
    if (err != nil) {
        return err
    }

    if (self.State > READY) {
        data = self.Cipher.Encrypt(data)
    }

    log.Debugf("send packet: %s", PacketToString(packet))

    err = self.SendLine(data)
    if (err != nil) {
        log.Error(err)
        return err
    }

    return nil;
}

func (self *ProtocolBase) SetPeer(protocol Gwisted.IProtocol, port int) {
    mu.Lock()
    defer mu.Unlock()

    log.Debugf("set peer for port %d", port)
    self.Outgoing[port] = protocol
}

func (self *ProtocolBase) RemovePeer(port int) {
    mu.Lock()
    defer mu.Unlock()

    log.Debugf("remove peer for port %d", port)
    if _, ok := self.Outgoing[port]; ok {
        delete(self.Outgoing, port)
    }
}

func (self *ProtocolBase) setTimeout() {
    self.timeoutTimer = time.AfterFunc(time.Millisecond * time.Duration(self.Timeout), self.timeoutConnection)
}

func (self *ProtocolBase) ResetTimeout() {
    self.timeoutTimer.Reset(time.Millisecond * time.Duration(self.Timeout))
}

func (self *ProtocolBase) LineReceived(data []byte) {
    log.Debugf("line received: %s", hex.EncodeToString(data))
    packet, err := self.parse(data)
    if (err != nil) {
        log.Error("packet parse error")
        self.HandleError()
    }

    switch self.State {
    case READY:
        self.LineReceivedOnReady(packet)
    case RUNNING:
        self.LineReceivedOnRunning(packet)
    default:
        log.Errorf("Unknown state: %d", self.State)
        self.HandleError()
    }
}

func (self *ProtocolBase) HandleError() {
    self.State = ERROR
    if (self.Transport != nil) {
        self.Transport.LoseConnection()
    } else {
        // clear all outgoings
    }
}

func (self *ProtocolBase) parse(data []byte) (*MPacket, error) {
    packet := &MPacket{}
    if (self.State > READY) {
        data = self.PeerCipher.Decrypt(data)
    }
    err := proto.Unmarshal(data, packet)
    return packet, err
}

func (self *ProtocolBase) timeoutConnection() {
    log.Error("connection timeout")
    self.HandleError()
}

func (self *ProtocolBase) LineReceivedOnReady(packet *MPacket) {
    if (self.LineReceivedOnReadyHandler == nil) {
        panic("LineReceivedOnReadyHandler is nil")
    }

    self.LineReceivedOnReadyHandler.LineReceivedOnReady(packet)
}

func (self *ProtocolBase) LineReceivedOnRunning(packet *MPacket) {
    if (self.LineReceivedOnRunningHandler == nil) {
        panic("LineReceivedOnRunningHandler is nil")
    }

    self.LineReceivedOnRunningHandler.LineReceivedOnRunning(packet)
}
