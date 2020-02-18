package Melkweg

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha256"
    "bytes"
    "fmt"
    "math/rand"
)

func DigestBytes(bytes []byte, length int) []byte {
    h := sha256.New()
    h.Write(bytes)
    digest := make([]byte, length)
    copy(digest, h.Sum(nil))
    return digest
}

func DigestString(str string, length int) []byte {
    bytes := []byte(str)
    return DigestBytes(bytes, length)
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

func NewAESCipher(iv []byte, key []byte) *AESCipher {
    block, err := aes.NewCipher(DigestBytes(key, 16))
    if (err != nil) {
        panic(err)
    }
    return &AESCipher {
        iv: iv,
        stream: cipher.NewCTR(block, iv),
    }
}

func (self *AESCipher) Encrypt(buffer []byte) []byte {
    if len(buffer) <= 0 {
        return make([]byte, 0)
    }
    encrypted := make([]byte, len(buffer))
    self.stream.XORKeyStream(encrypted, buffer)
    return encrypted
}

func (self *AESCipher) Decrypt(buffer []byte) []byte {
    if len(buffer) <= 0 {
        return make([]byte, 0)
    }
    decrypted := make([]byte, len(buffer))
    self.stream.XORKeyStream(decrypted, buffer)
    return decrypted
}

type AESHMacCipher struct {
    iv []byte
    key []byte
    stream cipher.Stream
}

func NewAESHMacCipher(iv []byte, key []byte) *AESHMacCipher {
    block, err := aes.NewCipher(DigestBytes(key, 16))
    if (err != nil) {
        panic(err)
    }
    return &AESHMacCipher {
        iv: iv,
        key: key,
        stream: cipher.NewCTR(block, iv),
    }
}

func (self *AESHMacCipher) Encrypt(buffer []byte) []byte {
    if len(buffer) <= 0 {
        return make([]byte, 0)
    }
    encrypted := make([]byte, len(buffer))
    self.stream.XORKeyStream(encrypted, buffer)
    tag := DigestBytes(append(append(self.key, buffer...), self.key...), 16)
    return append(encrypted, tag...)
}

func (self *AESHMacCipher) Decrypt(buffer []byte) []byte {
    if len(buffer) <= 0 {
        return make([]byte, 0)
    }
    length := len(buffer)
    encrypted := buffer[:length - 16]
    tag := buffer[length - 16:]

    decrypted := make([]byte, len(encrypted))
    self.stream.XORKeyStream(decrypted, encrypted)

    tag2 := DigestBytes(append(append(self.key, decrypted...), self.key...), 16)

    if (bytes.Compare(tag, tag2) != 0) {
        panic("hmac check error")
    }

    return decrypted
}

func CipherFactory(name string, iv []byte, key []byte) ICipher {
    if (name == "AES") {
        return NewAESCipher(iv, key)
    } else if (name == "AESHMac") {
        return NewAESHMacCipher(iv, key)
    } else {
        panic("unknown cipher name")
    }
}
