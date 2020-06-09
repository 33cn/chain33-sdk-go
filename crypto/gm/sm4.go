package gm

import (
	"github.com/tjfoc/gmsm/sm4"
	"log"
)

const (
	SM4PrivateKeyLength = 16
)

func GenetateSM4Key() []byte {
	return getRandBytes(SM4PrivateKeyLength)
}

func SM4Encrypt(key []byte, data []byte) []byte {
	c, err := sm4.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	cipher := make([]byte, 16)
	c.Encrypt(cipher, data)

	return cipher
}

func SM4Decrypt(key []byte, data []byte) []byte {
	c, err := sm4.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	cipher := make([]byte, 16)
	c.Decrypt(cipher, data)

	return cipher
}