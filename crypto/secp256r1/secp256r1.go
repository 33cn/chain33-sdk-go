package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/33cn/chain33-sdk-go/crypto/hash"
	"math/big"
)

func GeneratePrivateKey() ([]byte, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("Failed generating ECDSA key for [%v]: [%s]", elliptic.P256(), err)
	}
	privKeyBytes := make([]byte, privateKeyECDSALength)
	copy(privKeyBytes[:], SerializePrivateKey(privKey))
	return privKeyBytes, nil
}

func PubKeyFromPrivate(privKey []byte) []byte {
	_, pub := privKeyFromBytes(elliptic.P256(), privKey[:])
	pubSecp256r1 := make([]byte, publicKeyECDSALengthCompressed)
	copy(pubSecp256r1[:], SerializePublicKeyCompressed(pub))
	return pubSecp256r1
}

func Sign(msg []byte, privKey []byte) ([]byte, error) {
	priv, pub := privKeyFromBytes(elliptic.P256(), privKey[:])
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash.Sha256(msg))
	if err != nil {
		return nil, err
	}

	s = ToLowS(pub, s)
	ecdsaSigByte, err := MarshalECDSASignature(r, s)
	if err != nil {
		return nil, err
	}
	return ecdsaSigByte, nil
}

func Verify(msg []byte, sigECDSA []byte, pubKey []byte) bool {
	pub, err := parsePubKeyCompressed(pubKey[0:publicKeyECDSALengthCompressed])
	if err != nil {
		return false
	}

	r, s, err := UnmarshalECDSASignature(sigECDSA)
	if err != nil {
		return false
	}

	lowS := IsLowS(s)
	if !lowS {
		return false
	}
	return ecdsa.Verify(pub, hash.Sha256(msg), r, s)
}

func PrivateFromByte(privKey []byte) *ecdsa.PrivateKey {
	priv, _ := privKeyFromBytes(elliptic.P256(), privKey[:])

	return priv
}

func PublicFromByte(pubKey []byte) (*ecdsa.PublicKey, error) {
	return parsePubKeyCompressed(pubKey)
}

func getRandBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, err := rand.Read(b)
	if err != nil {
		panic("Panic on a Crisis" + err.Error())
	}
	return b
}

func privKeyFromBytes(curve elliptic.Curve, pk []byte) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	x, y := curve.ScalarBaseMult(pk)

	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}

	return priv, &priv.PublicKey
}
