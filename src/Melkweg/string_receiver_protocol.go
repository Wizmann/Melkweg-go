package Melkweg

import (
    "bytes"
    "net"
)

type IntNStringReceiver struct {
    isStopped   bool
    isPaused    bool
    buffer      bytes.Buffer
    readerCh    chan []byte
    conn        ITransport

    strSize     int
    prefixSize  uint32
    parsePrefix func([]byte) uint32
    makePrefix  func(uint32) []byte
    maxLength   uint32
}

func (self *IntNStringReceiver) dataReceived(data []byte) {
    length := len(data)

    if (length + self.buffer.Len() > maxLength) {
        self.LengthLimitExceeded()
        return
    }


    self.buffer.Write(data)

    if (strSize < 0 && self.buffer.Len() > prefixSize) {
        prefixBytes := make([]byte, prefixSize)
        self.buffer.Read(prefixBytes)
        self.strSize = self.parsePrefix(prefixBytes)
    }

    if (self.buffer.Len() < strSize) {
        return
    }

    lineData := make([]byte, strSize)
    self.buffer.Read(lineData)

    self.LineReceived(lineData)
}

func (self *IntNStringReceiver) LineReceived(data []byte) {
    // pass
}

func (self *IntNStringReceiver) lengthLimitExceeded() {
    self.transport.LoseConnection()
}

func (self *IntNStringReceiver) SendString(data []byte) (err error) {
    if (len(data) + self.prefixSize > self.maxLength) {
        return errors.New("Try to send %d bytes whereas max size limit is %u",
                len(data) + self.refixSize, self.maxLength)
    }

    self.transport.Write(append(self.makePrefix(len(data)), data))
    return nil;
}
