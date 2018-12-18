package Gwisted

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "testing"
)

type FakeTransport struct {
    Buffer []byte
    IsConnectionLost bool
}

func (self *FakeTransport) Write(data []byte) error {
    self.Buffer = data
    return nil
}

func (self *FakeTransport) LoseConnection() {
    self.IsConnectionLost = true
}

type MyLineReceiver struct {
    IntNStringReceiver

    Line []byte
}

func NewMyLineReceiver() *MyLineReceiver {
    buffer := make([]byte, 0)
    r := &MyLineReceiver {
        IntNStringReceiver: IntNStringReceiver {
            isStopped: false,
            isPaused: false,
            buffer: bytes.NewBuffer(buffer),
            readerCh: make(chan []byte),
            transport: &FakeTransport{},

            strSize: 99999,
            prefixSize: 4,
            parsePrefix: binary.BigEndian.Uint32,
            makePrefix: binary.BigEndian.PutUint32,
            maxLength: 99999,
        },
    }
    r.LineReceivedHandler = r

    return r
}

func (self *MyLineReceiver) LineReceived(data []byte) {
    self.Line = data
}

func TestLineReceiverProtocol(t *testing.T) {
    r := NewMyLineReceiver()

    prefix := make([]byte, 4)
    binary.BigEndian.PutUint32(prefix, 6)
    line := []byte("arbeit")

    r.dataReceived(append(prefix, line...))

    fmt.Println(r.strSize)
    fmt.Println(len(r.Line))
    if (bytes.Compare(r.Line, line) != 0) {
        t.Error()
    }
}