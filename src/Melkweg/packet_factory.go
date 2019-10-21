package Melkweg

import (
    "fmt"
)

const (
    SYN     = 0
    DATA    = 1
    LIV     = 2
    RST     = 3
    FIN     = 4
    KILL    = 5
)

func NewSynPacket(iv []byte) *MPacket {
    return &MPacket {
        Iv: iv,
        Flags: SYN,
    }
}

func NewRstPacket(port int) *MPacket {
    return &MPacket {
        Port: uint32(port),
        Flags: RST,
    }
}

func NewKillPacket() *MPacket {
    return &MPacket {
        Flags: KILL,
    }
}

func NewDataPacket(port int, data []byte) *MPacket {
    return &MPacket {
        Flags: DATA,
        Port: uint32(port),
        Data: data,
    }
}

func NewFinPacket(port int) *MPacket {
    return &MPacket {
        Flags: FIN,
        Port: uint32(port),
    }
}

func NewLivPacket() *MPacket {
    return &MPacket {
        Flags: LIV,
    }
}

func PacketToString(packet *MPacket) string {
    switch packet.Flags {
    case SYN:
        return fmt.Sprintf("[Syn Packet]")
    case DATA:
        return fmt.Sprintf("[Data Packet] %d bytes", len(packet.Data))
    case LIV:
        return fmt.Sprintf("[Liv Packet]")
    case FIN:
        return fmt.Sprintf("[Fin Packet] on port %d", packet.Port)
    case RST:
        return fmt.Sprintf("[Rst Packet] on port %d", packet.Port)
    }
    return "[Packet]"
}
