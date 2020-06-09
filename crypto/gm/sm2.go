package gm

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
)

const (
	SM2PrivateKeyLength = 32
)

var	DefaultUID = []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
	                    0x13, 0x23, 0x33, 0x43, 0x53, 0x63, 0x73, 0x83}

func getRandBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, err := rand.Read(b)
	if err != nil {
		panic("Panic on a Crisis" + err.Error())
	}
	return b
}

func privKeyFromBytes(curve elliptic.Curve, pk []byte) (*sm2.PrivateKey, *sm2.PublicKey) {
	x, y := curve.ScalarBaseMult(pk)

	priv := &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}

	return priv, &priv.PublicKey
}

func parsePubKey(pubKeyStr []byte) (key *sm2.PublicKey) {
	return sm2.Decompress(pubKeyStr)
}

//SerializePrivateKey 私钥序列化
func serializePrivateKey(p *sm2.PrivateKey) []byte {
	b := make([]byte, 0, SM2PrivateKeyLength)
	return paddedAppend(SM2PrivateKeyLength, b, p.D.Bytes())
}

//SerializePublicKey 公钥序列化
func serializePublicKey(p *sm2.PublicKey) []byte {
	return sm2.Compress(p)
}

func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

func canonicalizeInt(val *big.Int) []byte {
	b := val.Bytes()
	if len(b) == 0 {
		b = []byte{0x00}
	}
	if b[0]&0x80 != 0 {
		paddedBytes := make([]byte, len(b)+1)
		copy(paddedBytes[1:], b)
		b = paddedBytes
	}
	return b
}

func serializeSignature(r, s *big.Int) []byte {
	rb := canonicalizeInt(r)
	sb := canonicalizeInt(s)

	length := 6 + len(rb) + len(sb)
	b := make([]byte, length)

	b[0] = 0x30
	b[1] = byte(length - 2)
	b[2] = 0x02
	b[3] = byte(len(rb))
	offset := copy(b[4:], rb) + 4
	b[offset] = 0x02
	b[offset+1] = byte(len(sb))
	copy(b[offset+2:], sb)

	return b
}

func deserializeSignature(sigStr []byte) (*big.Int, *big.Int, error) {
	sig, err := btcec.ParseDERSignature(sigStr, sm2.P256Sm2())
	if err != nil {
		return nil, nil, err
	}

	return sig.R, sig.S, nil
}

func GenetateKey() ([]byte, []byte) {
	privKeyBytes := [SM2PrivateKeyLength]byte{}

	for {
		key := getRandBytes(32)
		if bytes.Compare(key, sm2.P256Sm2().Params().N.Bytes()) >= 0 {
			continue
		}
		copy(privKeyBytes[:], key)
		break
	}
	priv, pub := privKeyFromBytes(sm2.P256Sm2(), privKeyBytes[:])

	return serializePrivateKey(priv), serializePublicKey(pub)
}

func SM2Sign(msg []byte, privateKey []byte, uid []byte) []byte {
	if uid == nil {
		uid = DefaultUID
	}

	priv, _ := privKeyFromBytes(sm2.P256Sm2(), privateKey)
	r, s, err := sm2.Sm2Sign(priv, msg, uid)
	if err != nil {
		return nil
	}

	return serializeSignature(r, s)
}

func SM2Verify(msg []byte, publicKey []byte, sig []byte, uid []byte) bool {
	if uid == nil {
		uid = DefaultUID
	}

	pub := parsePubKey(publicKey[:])
	r, s, err := deserializeSignature(sig)
	if err != nil {
		fmt.Errorf("unmarshal sign failed")
		return false
	}

	return sm2.Sm2Verify(pub, msg, uid, r, s)
}

func SM2Encrypt(publicKey []byte, data []byte) ([]byte, error) {
	pub := parsePubKey(publicKey[:])

	return sm2.Encrypt(pub, data)
}

func SM2Decrypt(privateKey []byte, data []byte) ([]byte, error) {
	priv, _ := privKeyFromBytes(sm2.P256Sm2(), privateKey)

	return sm2.Decrypt(priv, data)
}

func PubKeyFromPrivate(privKey []byte) []byte {
	_, pub := privKeyFromBytes(sm2.P256Sm2(), privKey)
	return serializePublicKey(pub)
}