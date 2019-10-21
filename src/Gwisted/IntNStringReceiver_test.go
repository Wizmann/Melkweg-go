package Gwisted

import (
    "bytes"
    "encoding/binary"
    "testing"
)

type MyLineReceiver struct {
    IntNStringReceiver

    Line []byte
}

func NewMyLineReceiver() *MyLineReceiver {
    buffer := make([]byte, 0)
    r := &MyLineReceiver {
        IntNStringReceiver: IntNStringReceiver {
            buffer: bytes.NewBuffer(buffer),
            strSize: -1,
            prefixSize: 4,
            parsePrefix: func(buffer []byte) int {
                return int(binary.BigEndian.Uint32(buffer))
            },
            makePrefix: func(buffer []byte, prefix int) {
                binary.BigEndian.PutUint32(buffer, uint32(prefix))
            },
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

    r.DataReceived(append(prefix, line...))

    if (bytes.Compare(r.Line, line) != 0) {
        t.Error()
    }

    r.Line = nil

    binary.BigEndian.PutUint32(prefix, 10)
    r.DataReceived(prefix)
    if (r.Line != nil) {
        t.Error();
    }

    r.DataReceived([]byte("hello"))
    if (r.Line != nil) {
        t.Error();
    }

    r.DataReceived(append([]byte("world"), prefix...))
    if (bytes.Compare(r.Line, []byte("helloworld")) != 0) {
        t.Error()
    }
}
