package Melkweg

import (
    "fmt"
)

const (
    UNKNOWN = 0
    DATA    = 1
    LIV     = 2
    RST     = 3
    FIN     = 4
    KILL    = 5
)

func NewSynPacket(iv []byte) *MPacket {
    return &MPacket {
        Iv: iv,
    }
}

func NewRstPacket(port uint32) *MPacket {
    return &MPacket {
        Port: port,
        Flags: RST,
    }
}

func NewKillPacket() *MPacket {
    return &MPacket {
        Flags: KILL,
    }
}

func NewDataPacket(port uint32, data []byte) *MPacket {
    return &MPacket {
        Flags: DATA,
        Port: port,
        Data: data,
    }
}

func NewFinPacket(port uint32) *MPacket {
    return &MPacket {
        Flags: FIN,
        Port: port,
    }
}

func NewLivPacket() *MPacket {
    return &MPacket {
        Flags: LIV,
    }
}

func PacketToString(packet *MPacket) string {
    switch packet.Flags {
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
