package Gwisted

import (
    "bytes"
    "encoding/binary"
    "testing"
)

type FakeTransport struct {
    Buffer []data
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
    return &IntNStringReceiver {
        isStopped: false,
        isPaused: false,
        buffer: NewBuffer(make([]byte, 100 * 1024)),
        conn: &FakeTransport,

        strSize: -1,
        prefixSize: 4,
        parsePrefix: binary.BigEndian.Uint32
        makePrefix: binary.BigEndian.PutUint32
        maxLength = 99999
    }
}

func (self *MyLineReceiver) LineReceived(data []byte) {
    self.Line = data
}

func TestLineReceiverProtocol(t *testing.T) {
    r := NewMyLineReceiver()

    prefix := make([]byte, 4)
    binary.BigEndian.PutUint32(prefix, 6)
    line := []byte("arbeit")

    r.dataReceived(append(prefix, line))

    if (bytes.Compare(r.Line, line) != 0) {
        t.Error()
    }
}
