package crypto

import (
	"errors"

	"github.com/33cn/chain33-sdk-go/crypto/eth"
	"github.com/33cn/chain33-sdk-go/types"
	ccrypto "github.com/33cn/chain33/common/crypto"
	ttypes "github.com/33cn/chain33/types"
)

const (
	BTC_ADDRESS = 0
	ETH_ADDRESS = 2
)

type Account struct {
	PrivateKey string
	PublicKey  string
	Address    string
	Type       int32
}

func NewAccount(addrType int32) (*Account, error) {
	account := &Account{Type: addrType}
	if addrType == BTC_ADDRESS {
		priv := GeneratePrivateKey()
		pub := PubKeyFromPrivate(priv)
		addr, err := PubKeyToAddress(pub)
		if err != nil {
			return nil, err
		}

		account.PrivateKey = types.ToHex(priv)
		account.PublicKey = types.ToHex(pub)
		account.Address = addr
	} else if addrType == ETH_ADDRESS {
		priv := eth.GeneratePrivateKey()
		pub := eth.PubKeyFromPrivate(priv)
		addr, err := eth.PubKeyToAddress(pub)
		if err != nil {
			return nil, err
		}

		account.PrivateKey = types.ToHex(priv)
		account.PublicKey = types.ToHex(pub)
		account.Address = addr
	} else {
		return nil, errors.New("address type not support")
	}
	return account, nil
}

func SignTx(tx *ttypes.Transaction, privKey string, addrType int32) error {
	privkey, err := types.FromHex(privKey)
	if err != nil {
		return err
	}
	cr, err := ccrypto.Load(SECP256K1, -1)
	if err != nil {
		return err
	}
	priv, err := cr.PrivKeyFromBytes(privkey)
	if err != nil {
		return err
	}
	ty := ttypes.EncodeSignID(ttypes.SECP256K1, addrType)
	tx.Sign(ty, priv)
	return nil
}
