package Gwisted

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    _ "net"
)

type ILineReceivedHandler interface {
    LineReceived(data []byte)
}

type IntNStringReceiver struct {
    Protocol
    
    buffer      *bytes.Buffer

    strSize     int
    prefixSize  int
    parsePrefix func([]byte) int
    makePrefix  func([]byte, int)
    maxLength   int

    LineReceivedHandler ILineReceivedHandler
}

func (self *IntNStringReceiver) DataReceived(data []byte) {
    if (len(data) + self.buffer.Len() > self.maxLength) {
        self.lengthLimitExceeded()
        return
    }

    self.buffer.Write(data)

    if (self.strSize >= self.maxLength && self.buffer.Len() >= self.prefixSize) {
        prefixBytes := make([]byte, self.prefixSize)
        self.buffer.Read(prefixBytes)
        self.strSize = self.parsePrefix(prefixBytes)

        if (self.strSize < 0 || self.strSize > self.maxLength) {
            self.lengthLimitExceeded()
            return;
        }
    }

    if (self.buffer.Len() < self.strSize) {
        return
    }

    lineData := make([]byte, self.strSize)
    self.buffer.Read(lineData)
    self.strSize = self.maxLength

    self.LineReceived(lineData)
}

func (self *IntNStringReceiver) LineReceived(data []byte) {
    if (self.LineReceivedHandler != nil) {
        self.LineReceivedHandler.LineReceived(data);
    } else {
        // pass
    }
}

func (self *IntNStringReceiver) lengthLimitExceeded() {
    self.Transport.LoseConnection()
}

func (self *IntNStringReceiver) SendString(str string) error {
    return self.SendLine([]byte(str))
}

func (self *IntNStringReceiver) SendLine(data []byte) error {
    if (len(data) + self.prefixSize > self.maxLength) {
        return errors.New(
                fmt.Sprintf("Try to send %d bytes whereas max size limit is %d",
                    len(data) + self.prefixSize, self.maxLength))
    }

    prefix := make([]byte, self.prefixSize)
    self.makePrefix(prefix, len(data))
    self.Transport.Write(append(prefix, data...))
    return nil;
}

type Int32StringReceiver struct {
    IntNStringReceiver
}

func NewInt32StringReceiver() *Int32StringReceiver {
    r := &Int32StringReceiver {
        IntNStringReceiver: IntNStringReceiver {
            buffer: bytes.NewBuffer([]byte("")),
            strSize: 99999,
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
    r.DataReceivedHandler = r
    return r
}
