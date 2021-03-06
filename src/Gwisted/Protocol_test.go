package Gwisted

import (
    "bytes"
    _ "fmt"
    _ "encoding/binary"
    "testing"
    "time"
    logging "Logging"
)

type MyServerProtocol struct {
    Protocol
}

func (self *MyServerProtocol) DataReceived(data []byte) {
    logging.Debug("server data received: %x", data)
    if (bytes.Compare(data, []byte("ping")) == 0) {
        self.Transport.Write([]byte("pong"))
    } else if (bytes.Compare(data, []byte("exit")) == 0) {
        self.Transport.Write([]byte("exit"))
        self.Transport.LoseConnection()
    }
}

type MyClientProtocol struct {
    Protocol
}

func (self *MyClientProtocol) ConnectionMade(factory IProtocolFactory) {
    logging.Debug("client connection made")
    self.Transport.Write([]byte("ping"))
    self.Factory = factory
}

func (self *MyClientProtocol) DataReceived(data []byte) {
    logging.Debug("client data received %x", data)
    if (bytes.Compare(data, []byte("pong")) == 0) {
        self.Transport.Write([]byte("exit"))
    } else if (bytes.Compare(data, []byte("exit")) == 0) {
        self.Transport.LoseConnection()
    }
}

func NewMyServerProtocol() *MyServerProtocol {
    p := &MyServerProtocol{}
    p.DataReceivedHandler = p
    return p
}

func NewMyClientProtocol() *MyClientProtocol{
    p := &MyClientProtocol{}
    p.ConnectionMadeHandler = p
    p.DataReceivedHandler = p
    return p
}

func TestPingPongProtocol(t *testing.T) {
    server := NewMyServerProtocol()
    client := NewMyClientProtocol()

    if (server == nil || client == nil) {
        t.Error()
    }

    Reactor.ListenTCP(
        11000,
        ProtocolFactoryForProtocol(func() IProtocol { return server }),
        50)

    time.Sleep(time.Millisecond * 10)

    Reactor.ConnectTCP(
        "", 
        11000, 
        ProtocolFactoryForProtocol(func() IProtocol { return client }),
        60)

    time.Sleep(time.Millisecond * 1000)

    for i := 0; i <= 10; i++ {
        if (server.connected == 0 && client.connected == 0) {
            break;
        }
        if (i == 10) {
            t.Error();
        }
        time.Sleep(time.Millisecond * 500)
    }
}

type MyLineServerProtocol struct {
    *Int32StringReceiver
}

func (self *MyLineServerProtocol) LineReceived(data []byte) {
    logging.Debug("server data received: %x", data)
    if (bytes.Compare(data, []byte("ping")) == 0) {
        self.SendLine([]byte("pong"))
    } else if (bytes.Compare(data, []byte("exit")) == 0) {
        self.SendLine([]byte("exit"))
        self.Transport.LoseConnection()
    }
}

type MyLineClientProtocol struct {
    *Int32StringReceiver
}

func (self *MyLineClientProtocol) ConnectionMade(factory IProtocolFactory) {
    self.Factory = factory

    logging.Debug("client connection made")
    self.SendLine([]byte("ping"))
}

func (self *MyLineClientProtocol) LineReceived(data []byte) {
    logging.Debug("client data received: %x", data)
    if (bytes.Compare(data, []byte("pong")) == 0) {
        self.SendLine([]byte("exit"))
    } else if (bytes.Compare(data, []byte("exit")) == 0) {
        self.Transport.LoseConnection()
    }
}

func NewMyLineServerProtocol() *MyLineServerProtocol {
    p := &MyLineServerProtocol { NewInt32StringReceiver() }
    p.DataReceivedHandler = p
    p.LineReceivedHandler = p
    return p
}

func NewMyLineClientProtocol() *MyLineClientProtocol {
    p := &MyLineClientProtocol { NewInt32StringReceiver() }
    p.ConnectionMadeHandler = p
    p.LineReceivedHandler = p
    return p
}

func TestIntNLineProtocol(t *testing.T) {
    server := NewMyLineServerProtocol()
    client := NewMyLineClientProtocol()

    Reactor.ListenTCP(
        10000,
        ProtocolFactoryForProtocol(func() IProtocol { return server }),
        50)

    time.Sleep(time.Millisecond * 10)

    Reactor.ConnectTCP(
        "",
        10000,
        ProtocolFactoryForProtocol(func() IProtocol { return client }),
        60)

    time.Sleep(time.Millisecond * 1000)

    for i := 0; i <= 10; i++ {
        if (server.connected == 0 && client.connected == 0) {
            break
        }
        if (i == 10) {
            t.Error("server or client connection is still open")
        }
        time.Sleep(time.Millisecond * 500)
    }
}
