package Melkweg

import (
    logging "Logging"
    "time"
    "sync"
    "Gwisted"
    proto "github.com/golang/protobuf/proto"
)

type ProtocolState int;

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

type MelkwegProtocolBase struct {
    *Gwisted.Int32StringReceiver

    Key          []byte
    Iv           []byte
    Cipher       ICipher
    PeerCipher   ICipher
    State        ProtocolState
    Outgoing     map[int]Gwisted.IProtocol
    Timeout      int
    TimeoutTimer *time.Timer
    config       *ProxyConfig
    mu           sync.Mutex

    LineReceivedOnRunningHandler  ILineReceivedOnRunningHandler
    LineReceivedOnReadyHandler    ILineReceivedOnReadyHandler
}

func NewMelkwegProtocolBase(config *ProxyConfig) *MelkwegProtocolBase {
    p := &MelkwegProtocolBase {
        Int32StringReceiver: Gwisted.NewInt32StringReceiver(),
        config: config,
        Key: []byte(config.GetKey()),
        Timeout: config.GetTimeout(),
        Iv: DigestBytes(Nonce(19)),
        State: READY,
        Outgoing: map[int]Gwisted.IProtocol{},
    }
    p.LineReceivedHandler = p
    p.Cipher = NewAESCipher(p.Iv, p.Key)
    p.ConnectionMadeHandler = p
    p.ConnectionLostHandler = p

    p.setTimeout()

    return p
}

func (self *MelkwegProtocolBase) Write(packet *MPacket) error {
    self.mu.Lock()
    defer self.mu.Unlock()
    data, err := proto.Marshal(packet)
    if (err != nil) {
        return err
    }

    if (self.State > READY) {
        data = self.Cipher.Encrypt(data)
    }

    logging.Verbose("send packet: %s", PacketToString(packet))

    err = self.SendLine(data)
    if (err != nil) {
        logging.Error(err.Error())
        return err
    }

    return nil;
}

func (self *MelkwegProtocolBase) SetPeer(protocol Gwisted.IProtocol, port int) {
    logging.Debug("set peer for port %d", port)
    self.Outgoing[port] = protocol
}

func (self *MelkwegProtocolBase) RemovePeer(port int) {
    logging.Debug("remove peer for port %d", port)
    if _, ok := self.Outgoing[port]; ok {
        delete(self.Outgoing, port)
    }
}

func (self *MelkwegProtocolBase) setTimeout() {
    self.TimeoutTimer = time.AfterFunc(time.Millisecond * time.Duration(self.Timeout), self.timeoutConnection)
}

func (self *MelkwegProtocolBase) ResetTimeout() {
    self.TimeoutTimer.Reset(time.Millisecond * time.Duration(self.Timeout))
}

func (self *MelkwegProtocolBase) LineReceived(data []byte) {
    logging.Verbose("line received: %d bytes", len(data))
    packet, err := self.parse(data)
    if (err != nil) {
        logging.Error("packet parse error")
        self.HandleError()
    }

    switch self.State {
    case READY:
        self.LineReceivedOnReady(packet)
    case RUNNING:
        self.LineReceivedOnRunning(packet)
    default:
        logging.Error("Unknown state: %d", self.State)
        self.HandleError()
    }
}

func (self *MelkwegProtocolBase) HandleError() {
    self.State = ERROR
    logging.Error("handle error")
    if (self.Transport != nil) {
        self.Transport.LoseConnection()
    } else {
        // clear all outgoings
    }
}

func (self *MelkwegProtocolBase) parse(data []byte) (*MPacket, error) {
    packet := &MPacket{}
    if (self.State > READY) {
        data = self.PeerCipher.Decrypt(data)
    }
    err := proto.Unmarshal(data, packet)
    return packet, err
}

func (self *MelkwegProtocolBase) timeoutConnection() {
    logging.Error("connection timeout")
    self.HandleError()
}

func (self *MelkwegProtocolBase) LineReceivedOnReady(packet *MPacket) {
    if (self.LineReceivedOnReadyHandler == nil) {
        panic("LineReceivedOnReadyHandler is nil")
    }

    self.LineReceivedOnReadyHandler.LineReceivedOnReady(packet)
}

func (self *MelkwegProtocolBase) LineReceivedOnRunning(packet *MPacket) {
    if (self.LineReceivedOnRunningHandler == nil) {
        panic("LineReceivedOnRunningHandler is nil")
    }

    self.LineReceivedOnRunningHandler.LineReceivedOnRunning(packet)
}
