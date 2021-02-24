package  crypto

import (
	"fmt"
	"github.com/33cn/chain33-sdk-go/crypto/ed25519"
	"github.com/33cn/chain33-sdk-go/crypto/gm"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAES(t *testing.T) {
	var text = "hello aes"
	var key = getRandBytes(16)

	cipherText, err := AESCBCPKCS7Encrypt(key, []byte(text))
	if err != nil {
		fmt.Println(err)
		return
	}

	cipher, err := AESCBCPKCS7Decrypt(key, cipherText)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Equal(t, text, string(cipher))
}

func TestSign(t *testing.T) {
	priv, _ := types.FromHex("0xc2b31057b8692a56c7dd18199df71c1d21b781c0b6858c52997c9dbf778e8550")
	msg := []byte("sign test")

	sig := Sign(msg, priv)

	fmt.Printf("sig = %x\n", sig)

}

func TestSM2(t *testing.T) {
	priv, pub := gm.GenerateKey()
	fmt.Println(priv)
	fmt.Println(pub)

	fmt.Println(types.ToHex(priv))
	fmt.Println(types.ToHex(pub))

	msg := []byte("sign test")

	sig, _ := gm.SM2Sign(msg, priv,nil)
	fmt.Printf("sig = %x\n", sig)

	result := gm.SM2Verify(pub, msg, nil, sig)
	assert.Equal(t, true, result)
}

func TestSM4(t *testing.T) {
	key := []byte{0x1, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	fmt.Printf("key = %v\n", key)
	data := []byte{0x1, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	fmt.Printf("data = %x\n", data)
	d0 := gm.SM4Encrypt(key, data)
	fmt.Printf("d0 = %x\n", d0)
	d1 := gm.SM4Decrypt(key, d0)
	fmt.Printf("d1 = %x\n", d1)

	assert.Equal(t, data, d1)
}

func TestAddress(t *testing.T) {
	priv, _ := types.FromHex("0xc2b31057b8692a56c7dd18199df71c1d21b781c0b6858c52997c9dbf778e8550")

	pub := PubKeyFromPrivate(priv)
	fmt.Println(types.ToHex(pub))

	addr, err := PubKeyToAddress(pub)
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
}

func TestKDF(t *testing.T) {
	keyf := KDF([]byte("kdf test"), 16)
	fmt.Println(types.ToHex(keyf))
	assert.Equal(t, 16, len(keyf))
}

func TestED25519(t *testing.T) {
	priv, pub, err := ed25519.GenerateKey()
	if err != nil {
		assert.Fail(t, err.Error())
	}
	fmt.Println(types.ToHex(pub))

	msg := []byte("sign test")

	sig := ed25519.Sign(priv, msg)
	fmt.Printf("sig = %x\n", sig)

	result := ed25519.Verify(pub, msg, sig)
	assert.Equal(t, true, result)
}
