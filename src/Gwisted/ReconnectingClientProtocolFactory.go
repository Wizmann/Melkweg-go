package Gwisted

import (
    "net"
    "time"
)

type ReconnectingClientProtocolFactory struct {
    *ProtocolFactory

    initialDelay int
    delay        int
    maxDelay     int
    factor       float64
    maxRetries   int

    retries      int
}

func NewReconnectingClientProtocolFactory(
        protocolBuilder func(tcp *net.TCPConn) IProtocol,
        initialDelay int,
        delay int,
        maxDelay int,
        factor float64,
        maxRetries int) *ReconnectingClientProtocolFactory {
    f := &ReconnectingClientProtocolFactory{
        ProtocolFactory: NewProtocolFactory(protocolBuilder),
        delay: delay,
        initialDelay: initialDelay,
        maxDelay: maxDelay,
        factor: factor,
        maxRetries: maxRetries,
        retries: 0,
    }
    f.ClientConnectionLostHandler = f
    f.ClientConnectionFailedHandler = f
    return f
}

func NewReconnectingClientProtocolFactoryForProtocol(
        protocolCtor func() IProtocol,
        initialDelay int,
        delay int,
        maxDelay int,
        factor float64,
        maxRetries int) *ReconnectingClientProtocolFactory {
    f := &ReconnectingClientProtocolFactory{
        ProtocolFactory: ProtocolFactoryForProtocol(protocolCtor),
        delay: delay,
        initialDelay: initialDelay,
        maxDelay: maxDelay,
        factor: factor,
        maxRetries: maxRetries,
        retries: 0,
    }
    f.ClientConnectionLostHandler = f
    f.ClientConnectionFailedHandler = f
    return f
}

func (self *ReconnectingClientProtocolFactory) ClientConnectionLost(reason error) {
    self.retry()
}

func (self *ReconnectingClientProtocolFactory) ClientConnectionFailed(reason error) {
    self.retry()
}

func (self *ReconnectingClientProtocolFactory) retry() {
    self.retries += 1
    if (self.maxRetries > 0 && self.retries > self.maxRetries) {
        log.Errorf("abandon retry after %d retries", self.retries)
        return
    }
    delay := self.initialDelay + int(self.factor * float64(self.delay))
    if (delay > self.maxDelay) {
        delay = self.maxDelay
    }

    time.AfterFunc(time.Millisecond * time.Duration(delay), self.reconnect)
}

func (self *ReconnectingClientProtocolFactory) reconnect() {
    self.connector.Connect()
}