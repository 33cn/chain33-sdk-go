package types

import (
	"encoding/hex"
	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/golang/protobuf/proto"
)

//exec type
const (
	ExecErr  = 0
	ExecPack = 1
	ExecOk   = 2
)

//FromHex hex -> []byte
func FromHex(s string) ([]byte, error) {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
		if len(s)%2 == 1 {
			s = "0" + s
		}
		return hex.DecodeString(s)
	}
	return []byte{}, nil
}

//ToHex []byte -> hex
func ToHex(b []byte) string {
	hex := hex.EncodeToString(b)
	// Prefer output of "0x0" instead of "0x"
	if len(hex) == 0 {
		return ""
	}
	return hex
}

//ToHex []byte -> hex
func ToHexPrefix(b []byte) string {
	hex := hex.EncodeToString(b)
	// Prefer output of "0x0" instead of "0x"
	if len(hex) == 0 {
		return ""
	}
	return "0x" + hex
}

//Encode  编码
func Encode(data proto.Message) []byte {
	b, err := proto.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}

//Decode  解码
func Decode(data []byte, msg proto.Message) error {
	return proto.Unmarshal(data, msg)
}

// ECDH Calculate a shared secret using elliptic curve Diffie-Hellman
func ECDH(priv *secp256k1.PrivateKey, pub *secp256k1.PublicKey) []byte {
	ecKey := &secp256k1.PublicKey{}
	ecKey.X, ecKey.Y = secp256k1.S256().ScalarMult(pub.X, pub.Y, priv.D.Bytes())
	return ecKey.SerializeCompressed()
}