package  crypto

import (
	"fmt"
	"github.com/33cn/chain33-sdk-go/crypto/ed25519"
	"github.com/33cn/chain33-sdk-go/crypto/gm"
	"github.com/33cn/chain33-sdk-go/crypto/hash"
	"github.com/33cn/chain33-sdk-go/crypto/secp256r1"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAES(t *testing.T) {
	buf := getRandBytes(1e7)
	var key = getRandBytes(16)

	start := time.Now().UnixNano()
	cipherText, err := AESCBCPKCS7Encrypt(key, buf)
	if err != nil {
		assert.Error(t, err)
	}

	cipher, err := AESCBCPKCS7Decrypt(key, cipherText)
	if err != nil {
		assert.Error(t, err)
	}
	fmt.Println(time.Now().UnixNano() - start)

	assert.Equal(t, buf, cipher)
}

func BenchmarkAESEncrypt(b *testing.B) {
	var cipherText, cipher []byte
	var err error

	buf := getRandBytes(1024*1024*1024)
	var key = getRandBytes(16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cipherText, err = AESCBCPKCS7Encrypt(key, buf)
		if err != nil {
			assert.Error(b, err)
		}

		cipher, err = AESCBCPKCS7Decrypt(key, cipherText)
		if err != nil {
			assert.Error(b, err)
		}
	}
	assert.Equal(b, buf, cipher)
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
	data := getRandBytes(32)
	fmt.Printf("data = %x\n", data)
	d0 := gm.SM4Encrypt(key, data)
	fmt.Printf("d0 = %x\n", d0)
	d1 := gm.SM4Decrypt(key, d0)
	fmt.Printf("d1 = %x\n", d1)

	assert.Equal(t, data, d1)
}

func BenchmarkSM2Encrypt(b *testing.B) {
	var cipherText, cipher []byte
	var err error

	priv, pub := gm.GenerateKey()

	buf := getRandBytes(1024*1024*1024)
	fmt.Println(len(buf))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cipherText, err = gm.SM2Encrypt(pub, buf)
		if err != nil {
			assert.Error(b, err)
		}
		//fmt.Println(len(cipherText))
		//ioutil.WriteFile("cipherSm2.txt", cipherText, 666)

		cipher,err = gm.SM2Decrypt(priv, cipherText)
		if err != nil {
			assert.Error(b, err)
		}
	}
	assert.Equal(b, buf, cipher)
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
	keyf := hash.KDF([]byte("kdf test"), 16)
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

func TestSecp256r1(t *testing.T) {
	priv,err := secp256r1.GeneratePrivateKey()
	if err != nil {
		assert.Fail(t, err.Error())
	}
	pub := secp256r1.PubKeyFromPrivate(priv)
	fmt.Printf("priv = %x\n", priv)
	fmt.Printf("pub = %x\n", pub)

	msg := []byte("sign test")

	sig, err := secp256r1.Sign(msg, priv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	fmt.Printf("sig = %x\n", sig)

	result := secp256r1.Verify(msg, sig, pub)
	assert.Equal(t, true, result)
}