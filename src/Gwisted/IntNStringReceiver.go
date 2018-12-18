package Gwisted

// TODO: use int replace uint32

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

    strSize     uint32
    prefixSize  uint32
    parsePrefix func([]byte) uint32
    makePrefix  func([]byte, uint32)
    maxLength   uint32

    LineReceivedHandler ILineReceivedHandler
}

func (self *IntNStringReceiver) dataReceived(data []byte) {
    dataLength := uint32(len(data))
    bufferLength := uint32(self.buffer.Len())

    if (dataLength + bufferLength > self.maxLength) {
        self.lengthLimitExceeded()
        return
    }

    self.buffer.Write(data)
    bufferLength = uint32(self.buffer.Len())

    if (self.strSize >= self.maxLength && bufferLength >= self.prefixSize) {
        prefixBytes := make([]byte, self.prefixSize)
        self.buffer.Read(prefixBytes)
        self.strSize = self.parsePrefix(prefixBytes)
    }

    bufferLength = uint32(self.buffer.Len())

    if (bufferLength < self.strSize) {
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
    dataLength := uint32(len(data))
    if (dataLength + self.prefixSize > self.maxLength) {
        return errors.New(
                fmt.Sprintf("Try to send %d bytes whereas max size limit is %d",
                    uint32(len(data)) + self.prefixSize, self.maxLength))
    }

    prefix := make([]byte, self.prefixSize)
    self.makePrefix(prefix, uint32(len(data)))
    self.transport.Write(append(prefix, data...))
    return nil;
}
