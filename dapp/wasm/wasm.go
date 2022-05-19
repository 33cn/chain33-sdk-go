package wasm

import (
	"io/ioutil"
	"math/rand"
	"time"

	sdk "github.com/33cn/chain33-sdk-go"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
	ttypes "github.com/33cn/chain33/types"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func CreateWasmCreateTx(paraName, path, name string, privKey, cert, uid []byte) (*ttypes.Transaction, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	payload := &types.WasmAction{
		Ty: WasmActionCreate,
		Value: &types.WasmAction_Create{
			Create: &types.WasmCreate{
				Name: name,
				Code: code,
			},
		},
	}
	tx := &ttypes.Transaction{Execer: []byte(paraName + WasmX), Payload: types.Encode(payload), Fee: 1e5, Nonce: r.Int63(), To: crypto.GetExecAddress(paraName + WasmX)}
	tx, err = sdk.Sign(tx, privKey, crypto.SM2, uid)
	if err != nil {
		return nil, err
	}
	tx.Signature.Signature = crypto.EncodeCertToSignature(tx.Signature.Signature, cert, uid)
	return tx, nil
}

func CreateWasmUpdateTx(paraName, path, name string, privKey, cert, uid []byte) (*ttypes.Transaction, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	payload := &types.WasmAction{
		Ty: WasmActionUpdate,
		Value: &types.WasmAction_Update{
			Update: &types.WasmUpdate{
				Name: name,
				Code: code,
			},
		},
	}
	tx := &ttypes.Transaction{Execer: []byte(paraName + WasmX), Payload: types.Encode(payload), Fee: 1e5, Nonce: r.Int63(), To: crypto.GetExecAddress(paraName + WasmX)}
	tx, err = sdk.Sign(tx, privKey, crypto.SM2, uid)
	if err != nil {
		return nil, err
	}
	tx.Signature.Signature = crypto.EncodeCertToSignature(tx.Signature.Signature, cert, uid)
	return tx, nil
}

func CreateWasmCallTx(paraName, contract, method string, param []int64, env []string, privKey, cert, uid []byte) (*ttypes.Transaction, error) {
	payload := &types.WasmAction{
		Ty: WasmActionCall,
		Value: &types.WasmAction_Call{
			Call: &types.WasmCall{
				Contract:   contract,
				Method:     method,
				Parameters: param,
				Env:        env,
			},
		},
	}
	tx := &ttypes.Transaction{Execer: []byte(paraName + WasmX), Payload: types.Encode(payload), Fee: 1e5, Nonce: r.Int63(), To: crypto.GetExecAddress(paraName + WasmX)}
	var err error
	tx, err = sdk.Sign(tx, privKey, crypto.SM2, uid)
	if err != nil {
		return nil, err
	}
	tx.Signature.Signature = crypto.EncodeCertToSignature(tx.Signature.Signature, cert, uid)
	return tx, nil
}
