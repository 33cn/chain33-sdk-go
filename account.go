package sdk

import (
	"errors"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/crypto/gm"
)

type Account struct {
	PrivateKey  []byte
	PublicKey   []byte
	Address     string
	SignType    string
}

func NewAccount(signType string) (*Account, error) {
    if signType == "" {
    	signType = crypto.SECP256K1
	}

	account := Account{}
	account.SignType = signType
	if signType == crypto.SECP256K1 {
		account.PrivateKey = crypto.GeneratePrivateKey()
		account.PublicKey  = crypto.PubKeyFromPrivate(account.PrivateKey)

		addr, err := crypto.PubKeyToAddress(account.PublicKey)
		if err != nil {
			return nil, err
		}
		account.Address = addr
	} else if signType == crypto.SM2 {
		account.PrivateKey, account.PublicKey = gm.GenetateKey()
		addr, err := crypto.PubKeyToAddress(account.PublicKey)
		if err != nil {
			return nil, err
		}
		account.Address = addr
	} else if signType == crypto.ED25519 {
		// TODO
	} else {
		return nil, errors.New("sign type not support")
	}

	return &account, nil
}