package evm

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/33cn/chain33/common"
	ccrypto "github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	ttypes "github.com/33cn/chain33/types"
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
	"github.com/golang/protobuf/proto"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func CreateEvmContract(code []byte, note string, alias string, paraName string) (*ttypes.Transaction, error) {
	payload := &evmtypes.EVMContractAction{
		Code:         code,
		Note:         note,
		Alias:        alias,
		ContractAddr: crypto.GetExecAddress(paraName + EvmX),
	}
	tx := &ttypes.Transaction{Execer: []byte(paraName + EvmX), Payload: types.Encode(payload), Fee: EVM_FEE, Nonce: r.Int63(), To: crypto.GetExecAddress(paraName + EvmX)}
	return tx, nil
}

func CallEvmContract(param []byte, note string, amount int64, contractAddr string, paraName string) (*ttypes.Transaction, error) {
	payload := &evmtypes.EVMContractAction{
		Para:         param,
		Note:         note,
		Amount:       uint64(amount),
		ContractAddr: contractAddr,
	}
	tx := &ttypes.Transaction{Execer: []byte(paraName + EvmX), Payload: types.Encode(payload), Fee: EVM_FEE, Nonce: r.Int63(), To: crypto.GetExecAddress(paraName + EvmX)}
	return tx, nil
}

func EncodeParameter(abiStr, funcName string, params ...interface{}) ([]byte, error) {
	_, packedParameter, err := evmAbi.Pack(funcName, abiStr, false)
	return packedParameter, err
}

func GetContractAddr(deployer, hash, rpcLaddr string) (string, error) {
	params := &evmtypes.EvmCalcNewContractAddrReq{
		Caller: deployer,
		Txhash: hash,
	}
	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "evm.CalcNewContractAddr", params, &res)
	result, err := ctx.RunResult()
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func QueryContract(rpcLaddr, addr, abiStr, input, caller string) {
	methodName, packData, err := evmAbi.Pack(input, abiStr, true)
	if nil != err {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to do evmAbi.Pack")
		return
	}
	packStr := common.ToHex(packData)
	var req = evmtypes.EvmQueryReq{Address: addr, Input: packStr, Caller: caller}
	var resp evmtypes.EvmQueryResp

	query := sendQuery(rpcLaddr, "Query", &req, &resp)
	if !query {
		fmt.Println("Failed to send query")
		return

	}
	_, err = json.MarshalIndent(&resp, "", "  ")
	if err != nil {
		fmt.Println("MarshalIndent failed due to:", err.Error())
	}

	data, err := common.FromHex(resp.RawData)
	if nil != err {
		fmt.Println("common.FromHex failed due to:", err.Error())
	}

	outputs, err := evmAbi.Unpack(data, methodName, abiStr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "unpack evm return error", err)
	}

	for _, v := range outputs {
		fmt.Println(v.Value)
	}
}

func sendQuery(rpcAddr, funcName string, request ttypes.Message, result proto.Message) bool {
	params := rpctypes.Query4Jrpc{
		Execer:   "evm",
		FuncName: funcName,
		Payload:  ttypes.MustPBToJSON(request),
	}

	jsonrpc, err := jsonclient.NewJSONClient(rpcAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	err = jsonrpc.Call("Chain33.Query", params, result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	return true
}

func CreateNobalance(etx *ttypes.Transaction, fromAddressPriveteKey, withHoldPrivateKey, paraName string) (*ttypes.Transactions, error) {
	var noneExecer = "none"

	noneTx := &ttypes.Transaction{Execer: []byte(paraName + noneExecer), Payload: []byte("no-fee-transaction"), Nonce: rand.Int63()}
	noneTx.To = crypto.GetExecAddress(paraName + noneExecer)
	noneTx.Fee = EVM_FEE
	txs := []*ttypes.Transaction{noneTx}
	txs = append(txs, etx)

	group, err := ttypes.CreateTxGroup(txs, EVM_FEE)
	if err != nil {
		return nil, err
	}
	SignTx(group.Txs[0], withHoldPrivateKey)
	SignTx(group.Txs[1], fromAddressPriveteKey)

	return group, nil
}

func SignTx(tx *ttypes.Transaction, privKey string) error {
	privkey, err := types.FromHex(privKey)
	if err != nil {
		return err
	}
	cr, err := ccrypto.Load(crypto.SECP256K1, -1)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	priv, err := cr.PrivKeyFromBytes(privkey)
	if err != nil {
		return err
	}
	tx.Sign(ttypes.SECP256K1, priv)
	return nil
}
