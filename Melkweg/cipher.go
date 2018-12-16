package Cipher

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "fmt"
    "math/rand"
)

func DigestBytes(bytes []byte) []byte {
    h := md5.New()
    h.Write(bytes)
    return h.Sum(nil)
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

type ICipher interface {
    Encrypt(buffer []byte) []byte
    Decrypt(buffer []byte) []byte
}

type AESCipher struct {
    iv []byte
    stream cipher.Stream
}

func NewAESCipher(iv []byte, key string) *AESCipher {
    block, err := aes.NewCipher(DigestString(key))
    if (err != nil) {
        panic(err)
    }
    return &AESCipher {
        iv: iv,
        stream: cipher.NewCTR(block, iv),
    }
}

func (cipher *AESCipher) Encrypt(buffer []byte) []byte {
    encrypted := make([]byte, len(buffer))
    cipher.stream.XORKeyStream(encrypted, buffer)
    return encrypted
}

func (cipher *AESCipher) Decrypt(buffer []byte) []byte {
    decrypted := make([]byte, len(buffer))
    cipher.stream.XORKeyStream(decrypted, buffer)
    return decrypted
}
