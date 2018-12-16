package Cipher

import (
    "crypto/md5"
    "fmt"
    "math/rand"
)

func DigestBytes(bytes []byte) []byte {
    h := md5.New()
    return h.Sum(bytes)
}

func DigestString(str string) []byte {
    bytes := []byte(str)
    return DigestBytes(bytes)
}

func Nonce(length int) []byte {
    token := make([]byte, length)
    rand.Read(token)
    return token
}

func Hexlify(bytes []byte) string {
    return fmt.Sprintf("%x", bytes)
}
