package Gwisted

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "math"
    _ "net"
    logging "Logging"
)

type ILineReceivedHandler interface {
    LineReceived(data []byte)
}

type IntNStringReceiver struct {
    *Protocol

    buffer      *bytes.Buffer

    strSize     int
    prefixSize  int
    parsePrefix func([]byte) int
    makePrefix  func([]byte, int)
    maxLength   int

    LineReceivedHandler ILineReceivedHandler
}

func (self *IntNStringReceiver) DataReceived(data []byte) {
    logging.Debug("DataReceived: %d", len(data))

    if (len(data) + self.buffer.Len() > self.maxLength) {
        logging.Warning("String length error for IntNStringReceiver Protocol: %d", len(data) + self.buffer.Len())
        self.lengthLimitExceeded()
        return
    }

    self.buffer.Write(data)

    for {
        if (self.strSize == -1 && self.buffer.Len() >= self.prefixSize) {
            prefixBytes := make([]byte, self.prefixSize)
            self.buffer.Read(prefixBytes)
            self.strSize = self.parsePrefix(prefixBytes)

            if (self.strSize > self.maxLength) {
                logging.Warning("String length error for IntNStringReceiver Protocol: %d", self.strSize)
                self.lengthLimitExceeded()
                return;
            }
        }

        logging.Debug("buffer status, buffer length %d, strSize %d", self.buffer.Len(), self.strSize);

        if (self.strSize == -1 || self.buffer.Len() < self.strSize) {
            logging.Debug("buffer not ready, buffer length %d, strSize %d", self.buffer.Len(), self.strSize);
            return
        }

        lineData := make([]byte, self.strSize)
        self.buffer.Read(lineData)

        self.LineReceived(lineData)

        self.strSize = -1
    }
}

func (self *IntNStringReceiver) LineReceived(data []byte) {
    logging.Verbose("LineReceived: %x", data)
    if (self.LineReceivedHandler != nil) {
        self.LineReceivedHandler.LineReceived(data);
    } else {
        // pass
    }
}

func (self *IntNStringReceiver) lengthLimitExceeded() {
    logging.Fatal("length limit exceeded")
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
    logging.Verbose("length of data: %d", len(data))
    self.makePrefix(prefix, len(data))
    self.Transport.Write(append(prefix, data...))
    return nil
}

type Int32StringReceiver struct {
    IntNStringReceiver
}

func NewInt32StringReceiver() *Int32StringReceiver {
    r := &Int32StringReceiver {
        IntNStringReceiver: IntNStringReceiver {
            Protocol: NewProtocol(),
            buffer: bytes.NewBuffer([]byte("")),
            strSize: -1,
            prefixSize: 4,
            parsePrefix: func(buffer []byte) int {
                return int(binary.BigEndian.Uint32(buffer))
            },
            makePrefix: func(buffer []byte, prefix int) {
                binary.BigEndian.PutUint32(buffer, uint32(prefix))
            },
            maxLength: math.MaxInt32 / 2,
        },
    }
    r.DataReceivedHandler = r
    return r
}
