package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"io"
	"log"
)

func initIV() []byte {
	fiv := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, fiv); err != nil {
		panic(err.Error())
	}
	return fiv
}

func getKeyIV(token string, fiv []byte) ([]byte, []byte) {
	m1 := md5.New()
	m1.Write([]byte(token))
	m1.Write([]byte{0x01, 0x01, 0x01})
	key := m1.Sum(nil)

	m2 := md5.New()
	m2.Write([]byte(token))
	m2.Write(fiv)
	m2.Write([]byte{0x02, 0x02, 0x02})
	iv := m2.Sum(nil)

	return key, iv
}

func encryptSession(key, iv, p []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Crypto:", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Crypto:", err)
	}
	return aesgcm.Seal(nil, iv[:12], p, nil)
}

func decryptSession(key, iv, c []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Crypto:", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Crypto:", err)
	}

	plaintext, err := aesgcm.Open(nil, iv[:12], c, nil)
	return plaintext, err
}
