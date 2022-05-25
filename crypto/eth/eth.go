package eth

import (
	"crypto/ecdsa"
	"github.com/33cn/chain33-sdk-go/crypto"
	ccrypto "github.com/33cn/chain33/common/crypto"
	ttypes "github.com/33cn/chain33/types"

	"github.com/33cn/chain33-sdk-go/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

func GeneratePrivateKey() []byte {
	priv, _ := ethCrypto.GenerateKey()
	return ethCrypto.FromECDSA(priv)
}

func PubKeyFromPrivate(privKey []byte) []byte {
	priv, _ := ethCrypto.ToECDSA(privKey)
	pub := priv.PublicKey
	return ethCrypto.FromECDSAPub(&pub)
}

func Sign(msg []byte, privKey []byte) []byte {
	priv, _ := ethCrypto.ToECDSA(privKey)
	sig, _ := ethCrypto.Sign(msg, priv)
	return sig
}

func Validate(msg, pub, sig []byte) bool {
	newSig := sig[:len(sig)-1] // remove recovery id
	return ethCrypto.VerifySignature(pub, msg, newSig)
}

func PrivateFromByte(privKey []byte) *ecdsa.PrivateKey {
	priv, _ := ethCrypto.ToECDSA(privKey)
	return priv
}

func PublicFromByte(pubKey []byte) *ecdsa.PublicKey {
	pub, _ := ethCrypto.UnmarshalPubkey(pubKey)
	return pub
}

func PubKeyToAddress(pubKey []byte) string {
	pub, _ := ethCrypto.UnmarshalPubkey(pubKey)
	addr := ethCrypto.PubkeyToAddress(*pub)
	return types.ToHexPrefix(addr.Bytes())
}

func SignTx(tx *ttypes.Transaction, privKey string) error {
	privkey, err := types.FromHex(privKey)
	if err != nil {
		return err
	}
	cr, err := ccrypto.Load(crypto.SECP256K1, -1)
	if err != nil {
		return err
	}
	priv, err := cr.PrivKeyFromBytes(privkey)
	if err != nil {
		return err
	}
	ty := ttypes.EncodeSignID(ttypes.SECP256K1, 2)
	tx.Sign(ty, priv)
	return nil
}
