package Gwisted

type ITransport interface {
    Write(data []byte) error
    // WriteSequence(seq [][]byte) error
    LoseConnection()
    // GetPeer()
    // GetHost()
}
