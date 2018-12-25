package Melkweg

import (
    "bytes"
    "encoding/hex"
    "fmt"
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

func TestAESCipherWithServer(t *testing.T) {
    iv := []byte("0123456789abcdef")
    key := []byte("abcdef0123456789")
    cipher := NewAESCipher2(iv, key)

    encrypted1 := cipher.Encrypt([]byte("hello"))
    fmt.Println(hex.EncodeToString(encrypted1))

    if (hex.EncodeToString(encrypted1) != "8f07011604") {
        t.Errorf("%s not equal to %s", hex.EncodeToString(encrypted1), "6e89335a12")
    }
    // encrypted2 := cipher.Encrypt([]byte("hello"))

    // encrypted3 := cipher.Encrypt([]byte("world"))
}
