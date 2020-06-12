package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

const (
	AESPrivateKeyLength = 16
)

func pkcs7Padding(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func aesCBCEncrypt(key, s []byte) ([]byte, error) {
	return aesCBCEncryptWithRand(rand.Reader, key, s)
}

func aesCBCEncryptWithRand(prng io.Reader, key, s []byte) ([]byte, error) {
	if len(s)%aes.BlockSize != 0 {
		return nil, errors.New("Invalid plaintext. It must be a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(s))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(prng, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], s)

	return ciphertext, nil
}

func aesCBCDecrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(src) < aes.BlockSize {
		return nil, errors.New("Invalid ciphertext. It must be a multiple of the block size")
	}
	iv := src[:aes.BlockSize]
	src = src[aes.BlockSize:]

	if len(src)%aes.BlockSize != 0 {
		return nil, errors.New("Invalid ciphertext. It must be a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(src, src)

	return src, nil
}

func pkcs7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > aes.BlockSize || unpadding == 0 {
		return nil, errors.New("Invalid pkcs7 padding (unpadding > aes.BlockSize || unpadding == 0)")
	}

	pad := src[len(src)-unpadding:]
	for i := 0; i < unpadding; i++ {
		if pad[i] != byte(unpadding) {
			return nil, errors.New("Invalid pkcs7 padding (pad[i] != unpadding)")
		}
	}

	return src[:(length - unpadding)], nil
}

func AESCBCPKCS7Encrypt(key, src []byte) ([]byte, error) {
	tmp := pkcs7Padding(src)

	return aesCBCEncrypt(key, tmp)
}

func AESCBCPKCS7Decrypt(key, src []byte) ([]byte, error) {
	pt, err := aesCBCDecrypt(key, src)
	if err == nil {
		return pkcs7UnPadding(pt)
	}
	return nil, err
}

func GenetateAESKey() []byte {
	return getRandBytes(AESPrivateKeyLength)
}