package Melkweg

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

func TestAESCipher(t *testing.T) {
    iv := DigestBytes(Nonce(100))
    cipher1 := NewAESCipher(iv, "hello world")
    cipher2 := NewAESCipher(iv, "hello world")

    encrypted := cipher1.Encrypt([]byte("foo"))
    decrypted := cipher2.Decrypt(encrypted)

    if (bytes.Compare(decrypted, []byte("foo")) != 0) {
        t.Error()
    }

    encrypted = cipher1.Encrypt([]byte("bar"))
    decrypted = cipher2.Decrypt(encrypted)

    if (bytes.Compare(decrypted, []byte("bar")) != 0) {
        t.Error()
    }
}
