package evm

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	ttypes "github.com/33cn/chain33/types"
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmcommon "github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
	"github.com/golang/protobuf/proto"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func CreateEvmContract(code []byte, note string, alias string, paraName string, addressID int32, chainID int32) (*ttypes.Transaction, error) {
	payload := &evmtypes.EVMContractAction{
		Code:  code,
		Note:  note,
		Alias: alias,
	}
	payload.ContractAddr, _ = address.GetExecAddress(paraName+EvmX, addressID)
	tx := &ttypes.Transaction{Execer: []byte(paraName + EvmX), Payload: types.Encode(payload), Fee: 0, Nonce: r.Int63(), ChainID: chainID}
	tx.To, _ = address.GetExecAddress(paraName+EvmX, addressID)
	return tx, nil
}

func CallEvmContract(param []byte, note string, amount int64, contractAddr string, paraName string, addressID int32, chainID int32) (*ttypes.Transaction, error) {
	payload := &evmtypes.EVMContractAction{
		Para:         param,
		Note:         note,
		Amount:       uint64(amount),
		ContractAddr: contractAddr,
	}
	tx := &ttypes.Transaction{Execer: []byte(paraName + EvmX), Payload: types.Encode(payload), Fee: 0, Nonce: r.Int63(), ChainID: chainID}
	tx.To, _ = address.GetExecAddress(paraName+EvmX, addressID)
	return tx, nil
}

func EncodeParameter(abiStr, funcName string, params ...interface{}) ([]byte, error) {
	_, packedParameter, err := evmAbi.Pack(funcName, abiStr, false)
	return packedParameter, err
}

func LocalGetContractAddr(caller string, txhash []byte, addrType int32) string {
	InitAddrType(addrType)
	return evmcommon.NewContractAddress(*evmcommon.StringToAddress(caller), txhash).String()
}

func InitAddrType(addrType int32) {
	driver, err := address.LoadDriver(addrType, -1)
	if err != nil {
		panic(err)
	}
	evmcommon.InitEvmAddressTypeOnce(driver)
}

func GetContractAddr(deployer, hash, rpcLaddr string) string {
	params := &evmtypes.EvmCalcNewContractAddrReq{
		Caller: deployer,
		Txhash: hash,
	}
	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "evm.CalcNewContractAddr", params, &res)
	ctx.RunResult()
	return res
}

func GetProperFee(rpcLaddr string) int64 {
	var res rpctypes.ReplyProperFee
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.GetProperFee", nil, &res)
	ctx.RunResult()
	return res.ProperFee
}

func QueryEvmGas(rpcLaddr, txStr, caller string) (int64, error) {
	txInfo := &evmtypes.EstimateEVMGasReq{
		Tx:   txStr,
		From: caller,
	}

	var estGasResp evmtypes.EstimateEVMGasResp
	err := sendQuery(rpcLaddr, "EstimateGas", txInfo, &estGasResp)
	if err != nil {
		return 0, fmt.Errorf("gas cost estimate error: %s", err)
	}
	return int64(estGasResp.Gas), nil
}

func UpdateTxFee(tx *ttypes.Transaction, gas int64, fee int64) {
	realFee := int64(0)
	if fee > gas {
		realFee = fee
	} else {
		realFee = gas
	}
	tx.Fee = realFee + EVM_FEE
}

func QueryContract(rpcLaddr, addr, abiStr, input, caller string) ([]interface{}, error) {
	methodName, packData, err := evmAbi.Pack(input, abiStr, true)
	if err != nil {
		return nil, fmt.Errorf("Failed to do evmAbi.Pack")
	}
	packStr := common.ToHex(packData)
	var req = evmtypes.EvmQueryReq{Address: addr, Input: packStr, Caller: caller}
	var resp evmtypes.EvmQueryResp

	err = sendQuery(rpcLaddr, "Query", &req, &resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to send query: %s", err)

	}
	_, err = json.MarshalIndent(&resp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("MarshalIndent failed due to: %s", err)
	}

	data, err := common.FromHex(resp.RawData)
	if nil != err {
		return nil, fmt.Errorf("common.FromHex failed due to: %s", err)
	}

	outputs, err := evmAbi.Unpack(data, methodName, abiStr)
	if err != nil {
		return nil, fmt.Errorf("unpack evm error: %s", err)
	}

	ret := make([]interface{}, 0)
	for _, v := range outputs {
		ret = append(ret, v.Value)
	}
	return ret, nil
}

func sendQuery(rpcAddr, funcName string, request ttypes.Message, result proto.Message) error {
	params := rpctypes.Query4Jrpc{
		Execer:   "evm",
		FuncName: funcName,
		Payload:  ttypes.MustPBToJSON(request),
	}

	jsonrpc, err := jsonclient.NewJSONClient(rpcAddr)
	if err != nil {
		return err
	}

	err = jsonrpc.Call("Chain33.Query", params, result)
	if err != nil {
		return err
	}
	return nil
}

func CreateNobalance(etx *ttypes.Transaction, fromAddressPriveteKey, withHoldPrivateKey, paraName string, addressID int32, chainID int32) (*ttypes.Transactions, error) {
	var noneExecer = "none"
	noneTx := &ttypes.Transaction{Execer: []byte(paraName + noneExecer), Payload: []byte("no-fee-transaction"), Nonce: rand.Int63(), ChainID: chainID}
	noneTx.To, _ = address.GetExecAddress(paraName+noneExecer, addressID)
	noneTx.Fee = etx.Fee
	txs := []*ttypes.Transaction{noneTx}
	txs = append(txs, etx)

	group, err := ttypes.CreateTxGroup(txs, etx.Fee)
	if err != nil {
		return nil, err
	}
	err = crypto.SignTx(group.Txs[0], withHoldPrivateKey, addressID)
	if err != nil {
		return nil, err
	}
	err = crypto.SignTx(group.Txs[1], fromAddressPriveteKey, addressID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

