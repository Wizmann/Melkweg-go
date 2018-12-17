package Melkweg

import (
    "bytes"
    "testing"
)

func TestNewPacket(t *testing.T) {
    syn := NewSynPacket([]byte("hello world"))
    if (bytes.Compare(syn.Iv, []byte("hello world")) != 0) {
        t.Error()
    }
}
