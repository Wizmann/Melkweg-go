package Melkweg

const (
    UNKNOWN uint32 = iota
    DATA
    LIV
    RST
    FIN
    KILL
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
