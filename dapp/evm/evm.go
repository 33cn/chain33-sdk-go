package evm

import (
	"encoding/json"
	"fmt"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/rpc/jsonclient"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"os"
	"time"

	sdk "github.com/33cn/chain33-sdk-go"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
	rpctypes "github.com/33cn/chain33/rpc/types"
	ttypes "github.com/33cn/chain33/types"
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func CreateEvmContract(code []byte, note string, alias string, privateKey string, paraName string, gas int64) (*ttypes.Transaction, error) {
	payload := &evmtypes.EVMContractAction{
		Code:         code,
		Note:         note,
		Alias:        alias,
		ContractAddr: crypto.GetExecAddress(paraName + EvmX),
	}
	var fee int64
	if gas < EVM_FEE {
		fee = EVM_FEE
	} else {
		fee = gas + 100000
	}
	tx := &ttypes.Transaction{Execer: []byte(paraName + EvmX), Payload: types.Encode(payload), Fee: fee, Nonce: r.Int63(), To: crypto.GetExecAddress(paraName + EvmX)}
	privByte, _ := types.FromHex(privateKey)
	sdk.Sign(tx, privByte, crypto.SECP256K1, nil)
	return tx, nil
}

func CallEvmContract(param []byte, note string, amount int64, contractAddr string, privateKey string, paraName string, gas int64) (*ttypes.Transaction, error) {
	payload := &evmtypes.EVMContractAction{
		Para:         param,
		Note:         note,
		Amount:       uint64(amount),
		ContractAddr: contractAddr,
	}
	var fee int64
	if gas < EVM_FEE {
		fee = EVM_FEE
	} else {
		fee = gas + 100000
	}
	tx := &ttypes.Transaction{Execer: []byte(paraName + EvmX), Payload: types.Encode(payload), Fee: fee, Nonce: r.Int63(), To: crypto.GetExecAddress(paraName + EvmX)}
	privByte, _ := types.FromHex(privateKey)
	sdk.Sign(tx, privByte, crypto.SECP256K1, nil)
	return tx, nil
}

func EncodeParameter(abiStr, funcName string, params ...interface{}) ([]byte, error) {
	_, packedParameter, err := evmAbi.Pack(funcName, abiStr, false)
	return packedParameter, err
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
