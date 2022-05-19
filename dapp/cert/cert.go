package cert

import (
	sdk "github.com/33cn/chain33-sdk-go"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
	ttypes "github.com/33cn/chain33/types"
	"math/rand"
	"time"
)

func CreateCertNormalTx(paraName string, privateKey []byte, cert []byte, uid []byte, key string, value []byte) (*ttypes.Transaction, error) {
	payload := &types.CertAction{
		Value: &types.CertAction_Normal{
			&types.CertNormal{
				Key:   key,
				Value: value,
			},
		},
		Ty: CertActionNormal,
	}

	var tx *ttypes.Transaction
	if paraName == "" {
		tx = &ttypes.Transaction{Execer: []byte(CertX), Payload: types.Encode(payload), Fee: 1e5, Nonce: rand.Int63n(time.Now().UnixNano()), To: crypto.GetExecAddress(CertX)}
	} else {
		tx = &ttypes.Transaction{Execer: []byte(paraName + CertX), Payload: types.Encode(payload), Fee: 1e5, Nonce: rand.Int63n(time.Now().UnixNano()), To: crypto.GetExecAddress(paraName + CertX)}
	}

	var err error
	tx, err = sdk.Sign(tx, privateKey, crypto.SM2, uid)
	if err != nil {
		return nil, err
	}

	tx.Signature.Signature = crypto.EncodeCertToSignature(tx.Signature.Signature, cert, uid)

	return tx, nil
}
