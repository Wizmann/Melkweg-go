package Gwisted

import (
    "bytes"
    "testing"
    "time"
    logging "Logging"
)

type ServerProtocol struct {
    *Int32StringReceiver
}

func NewServerProtocol() *ServerProtocol {
    p := &ServerProtocol {
        Int32StringReceiver: NewInt32StringReceiver(),
    }
    p.ConnectionMadeHandler = p
    p.LineReceivedHandler = p

    return p
}

func (self *ServerProtocol) ConnectionMade(factory IProtocolFactory) {
    self.Factory = factory

    logging.Debug("server connection made")
    self.SendLine([]byte("ping"))
}

func (self *ServerProtocol) LineReceived(data []byte) {
    logging.Verbose("client data received: %x", data)
    if (bytes.Compare(data, []byte("pong")) == 0) {
        self.SendLine([]byte("ping"))
    } else if (bytes.Compare(data, []byte("exit")) == 0) {
        self.Transport.LoseConnection()
    }
}

type ClientProtocol struct {
    *Int32StringReceiver

    loop int
}

func NewClientProtocol() *ClientProtocol {
    p := &ClientProtocol {
        Int32StringReceiver: NewInt32StringReceiver(),
        loop: 100,
    }
    p.LineReceivedHandler = p

    return p
}

func (self *ClientProtocol) LineReceived(data []byte) {
    logging.Verbose("client data received: %x", data)
    if (bytes.Compare(data, []byte("ping")) == 0) {
        if (self.loop == 0) {
            self.SendLine([]byte("exit"))
        } else {
            self.SendLine([]byte("pong"))
            self.loop -= 1
        }
    } else {
        self.Transport.LoseConnection()
    }
    logging.Warning("loop: %d", self.loop)
}


func TestLineBasedPingPong(t *testing.T) {
    var clientProtocol = NewClientProtocol()
    var serverProtocol = NewServerProtocol()

    var clientProtocolCreator = func() IProtocol {
        return clientProtocol
    }

    var serverProtocolCreator = func() IProtocol {
        return serverProtocol
    }

    Reactor.ListenTCP(
        11111, ProtocolFactoryForProtocol(serverProtocolCreator), 5)

    time.Sleep(time.Duration(100) * time.Millisecond)

    Reactor.ConnectTCP(
        "localhost", 11111, ProtocolFactoryForProtocol(clientProtocolCreator), 1000)

    for i := 0; i < 1000; i++ {
        if (clientProtocol.loop == 0) {
            break
        }
        time.Sleep(time.Duration(100) * time.Millisecond)
    }

    if (clientProtocol.loop != 0) {
        t.Error()
    }
}
