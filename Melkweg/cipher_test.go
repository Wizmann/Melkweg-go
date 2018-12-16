package Cipher

import (
    "bytes"
    "testing"
)

func TestDigest(t *testing.T) {
    str := "hello world"
    bin := []byte(str)

    if (bytes.Compare(DigestString(str), DigestBytes(bin)) != 0) {
        t.Error()
    }
}
