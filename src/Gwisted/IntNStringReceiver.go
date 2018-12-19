package Gwisted

import (
    "bytes"
    "errors"
    "fmt"
    _ "net"
)

type ILineReceivedHandler interface {
    LineReceived(data []byte)
}

type IntNStringReceiver struct {
    isStopped   bool
    isPaused    bool
    buffer      *bytes.Buffer
    readerCh    chan []byte
    transport   ITransport

    strSize     int
    prefixSize  int
    parsePrefix func([]byte) int
    makePrefix  func([]byte, int)
    maxLength   int

    LineReceivedHandler ILineReceivedHandler
}

func (self *IntNStringReceiver) dataReceived(data []byte) {
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
    self.transport.LoseConnection()
}

func (self *IntNStringReceiver) SendString(data []byte) (err error) {
    if (len(data) + self.prefixSize > self.maxLength) {
        return errors.New(
                fmt.Sprintf("Try to send %d bytes whereas max size limit is %d",
                    len(data) + self.prefixSize, self.maxLength))
    }

    prefix := make([]byte, self.prefixSize)
    self.makePrefix(prefix, len(data))
    self.transport.Write(append(prefix, data...))
    return nil;
}
