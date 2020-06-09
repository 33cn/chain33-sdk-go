package crypto

import (
	"bytes"
	"crypto/rand"
	secp256k1 "github.com/btcsuite/btcd/btcec"
)

var (
	SECP256K1 = "secp256k1"
	SM2       = "sm2"
	ED25519   = "ed25519" //TODO
)

func GeneratePrivateKey() []byte {
	privKeyBytes := make([]byte, 32)
	for {
		key := getRandBytes(32)
		if bytes.Compare(key, secp256k1.S256().Params().N.Bytes()) >= 0 {
			continue
		}
		copy(privKeyBytes[:], key)
		break
	}

	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKeyBytes[:])
	copy(privKeyBytes[:], priv.Serialize())
	return privKeyBytes
}

func PubKeyFromPrivate(privKey []byte) []byte {
	_, pub := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	pubSecp256k1 := make([]byte, 33)
	copy(pubSecp256k1[:], pub.SerializeCompressed())
	return pubSecp256k1
}

func Sign(msg []byte, privKey []byte) []byte {
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	sig, err := priv.Sign(Sha256(msg))
	if err != nil {
		panic("Error signing secp256k1" + err.Error())
	}
	return sig.Serialize()
}

func PrivateFromByte(privKey []byte) *secp256k1.PrivateKey {
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	return priv
}

func PublicFromByte(pubKey []byte) *secp256k1.PublicKey {
	pub, _ := secp256k1.ParsePubKey(pubKey, secp256k1.S256())
	return pub
}

func getRandBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, err := rand.Read(b)
	if err != nil {
		panic("Panic on a Crisis" + err.Error())
	}
	return b
}