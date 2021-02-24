package sdk

import (
	"errors"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/crypto/ed25519"
	"github.com/33cn/chain33-sdk-go/crypto/gm"
	"github.com/33cn/chain33-sdk-go/types"
	log "github.com/inconshreveable/log15"
)

type Account struct {
	PrivateKey  []byte
	PublicKey   []byte
	Address     string
	SignType    string
}

var rlog = log.New("module", "chain33 adk")

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
			rlog.Error("NewAccount.PubKeyToAddress", "error", err.Error())
			return nil, err
		}
		account.Address = addr
	} else if signType == crypto.SM2 {
		account.PrivateKey, account.PublicKey = gm.GenerateKey()
		addr, err := crypto.PubKeyToAddress(account.PublicKey)
		if err != nil {
			rlog.Error("NewAccount.PubKeyToAddress", "error", err.Error())
			return nil, err
		}
		account.Address = addr
	} else if signType == crypto.ED25519 {
		priv, pub, err := ed25519.GenerateKey()
		if err != nil {
			rlog.Error("NewAccount.GenerateKey", "error", err.Error())
			return nil, err
		}
		copy(account.PrivateKey, priv)
		copy(account.PublicKey, pub)

		addr, err := crypto.PubKeyToAddress(account.PublicKey)
		if err != nil {
			rlog.Error("NewAccount.PubKeyToAddress", "error", err.Error())
			return nil, err
		}
		account.Address = addr
	} else {
		rlog.Error("sign type not support")
		return nil, errors.New("sign type not support")
	}

	return &account, nil
}

func NewAccountFromLocal(signType string, filePath string) (*Account, error) {
	if signType == "" {
		signType = crypto.SECP256K1
	}

	account := Account{}
	account.SignType = signType

	if signType == crypto.SECP256K1 {
		//TODO
		return nil, errors.New("not support")
	} else if signType == crypto.SM2 {
		content, err := types.ReadFile(filePath)
		if err != nil {
			rlog.Error("GetKeyByte.read key file failed.", "file", filePath, "error", err.Error())
			return nil, err
		}

		keyBytes,err := types.FromHex(string(content))
		if err != nil {
			rlog.Error("GetKeyByte.FromHex.", "error", err.Error())
			return nil, err
		}

		if len(keyBytes) != gm.SM2PrivateKeyLength {
			rlog.Error("GetKeyByte.private key length error", "len", len(keyBytes), "expect", gm.SM2PrivateKeyLength)
			return nil, errors.New("private key length error")
		}
		account.PrivateKey = keyBytes

		account.PublicKey = gm.PubKeyFromPrivate(keyBytes)
		addr, err := crypto.PubKeyToAddress(account.PublicKey)
		if err != nil {
			rlog.Error("NewAccount.PubKeyToAddress", "error", err.Error())
			return nil, err
		}
		account.Address = addr
	} else if signType == crypto.ED25519 {
		return nil, errors.New("not support")
	} else {
		return nil, errors.New("sign type not support")
	}

	return &account, nil
}
