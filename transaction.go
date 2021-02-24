package sdk

import (
	"errors"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/crypto/ed25519"
	"github.com/33cn/chain33-sdk-go/crypto/gm"
	. "github.com/33cn/chain33-sdk-go/types"
)

func Sign(tx *Transaction, privateKey []byte, signType string, uid []byte) (*Transaction, error) {
	if signType == "" {
		signType = crypto.SECP256K1
	}

	tx.Signature = nil
	if signType == crypto.SECP256K1 {
		pub := crypto.PubKeyFromPrivate(privateKey)

		data := Encode(tx)
		signature := crypto.Sign(data, privateKey)
		tx.Signature = &Signature{
			Ty:        1,
			Pubkey:    pub,
			Signature: signature,
		}
	} else if signType == crypto.SM2 {
		pub := gm.PubKeyFromPrivate(privateKey)

		data := Encode(tx)
		signature, err := gm.SM2Sign(data, privateKey, uid)
		if err != nil {
			return nil, err
		}
		tx.Signature = &Signature{
			Ty:        258,
			Pubkey:    pub,
			Signature: signature,
		}
	} else if signType == crypto.ED25519 {
		pub := ed25519.MakePublicKey(privateKey)
		data := Encode(tx)
		signature := ed25519.Sign(data, privateKey)
		tx.Signature = &Signature{
			Ty:                   2,
			Pubkey:               pub,
			Signature:            signature,
		}
	} else {
		return nil, errors.New("sign type not support")
	}

	return tx, nil
}

func cloneTx(tx *Transaction) *Transaction {
	copytx := &Transaction{}
	copytx.Execer = tx.Execer
	copytx.Payload = tx.Payload
	copytx.Signature = tx.Signature
	copytx.Fee = tx.Fee
	copytx.Expire = tx.Expire
	copytx.Nonce = tx.Nonce
	copytx.To = tx.To
	copytx.GroupCount = tx.GroupCount
	copytx.Header = tx.Header
	copytx.Next = tx.Next
	return copytx
}

func Hash(tx *Transaction) []byte {
	copytx := cloneTx(tx)
	copytx.Signature = nil
	copytx.Header = nil
	data := Encode(copytx)
	return crypto.Sha256(data)
}