package Melkweg

import (
    "fmt"
    "encoding/hex"
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
    cipher1 := NewAESCipher(iv, DigestString("hello world"))
    cipher2 := NewAESCipher(iv, DigestString("hello world"))

    encrypted := cipher1.Encrypt([]byte())
    decrypted := cipher2.Decrypt(encrypted)

    encrypted = cipher1.Encrypt([]byte("foo"))
    decrypted = cipher2.Decrypt(encrypted)

    if (bytes.Compare(decrypted, []byte("foo")) != 0) {
        t.Error()
    }

    encrypted = cipher1.Encrypt([]byte("bar"))
    decrypted = cipher2.Decrypt(encrypted)

    if (bytes.Compare(decrypted, []byte("bar")) != 0) {
        t.Error()
    }
}

func TestAESCipherCorrectness(t *testing.T) {
    iv := []byte("0123456789abcdef")
    key := []byte("abcdef0123456789")
    cipher := NewAESCipher(iv, key)

    encrypted1 := cipher.Encrypt([]byte("hello"))
    fmt.Println(hex.EncodeToString(encrypted1))

    encrypted2 := "ae588a2b2b"

    if (hex.EncodeToString(encrypted1) != encrypted2) {
        t.Errorf("%s not equal to %s", hex.EncodeToString(encrypted1), encrypted2)
    }
}
