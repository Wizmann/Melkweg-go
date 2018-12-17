package Melkweg

type ProtocolState int;

const (
    READY   = ProtocolState(1)
    RUNNING = ProtocolState(2)
    DONE    = ProtocolState(3)
    ERROR   = ProtocolState(4)
)

type Protocol struct {
    key string
    iv []byte
    cipher ICipher
    state ProtocolState
}
