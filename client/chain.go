package client

import (
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	commandtypes "github.com/33cn/chain33/system/dapp/commands/types"
	"github.com/33cn/chain33/types"
)

func (client *JSONClient) GetLastHeader() (*rpctypes.Header, error) {
	var res rpctypes.Header
	ctx := jsonclient.NewRPCCtx(client.url, "Chain33.GetLastHeader", nil, &res)
	_, err := ctx.RunResult()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (client *JSONClient) GetBlocks(start, end int64, isDetail bool) (*rpctypes.BlockDetails, error) {
	params := rpctypes.BlockParam{
		Start:    start,
		End:      end,
		Isdetail: isDetail,
	}
	var res rpctypes.BlockDetails
	ctx := jsonclient.NewRPCCtx(client.url, "Chain33.GetBlocks", params, &res)
	ctx.SetResultCb(parseBlockDetail)
	_, err := ctx.RunResult()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func parseBlockDetail(res interface{}) (interface{}, error) {
	var result commandtypes.BlockDetailsResult
	for _, vItem := range res.(*rpctypes.BlockDetails).Items {
		b := &commandtypes.BlockResult{
			Version:    vItem.Block.Version,
			ParentHash: vItem.Block.ParentHash,
			TxHash:     vItem.Block.TxHash,
			StateHash:  vItem.Block.StateHash,
			Height:     vItem.Block.Height,
			BlockTime:  vItem.Block.BlockTime,
		}
		for _, vTx := range vItem.Block.Txs {
			b.Txs = append(b.Txs, commandtypes.DecodeTransaction(vTx))
		}
		bd := &commandtypes.BlockDetailResult{Block: b, Receipts: vItem.Receipts}
		result.Items = append(result.Items, bd)
	}
	return result, nil
}

func (client *JSONClient) GetBlockHashByHeight(height int64) (string, error) {
	params := types.ReqInt{
		Height: height,
	}
	var res rpctypes.ReplyHash
	ctx := jsonclient.NewRPCCtx(client.url, "Chain33.GetBlockHash", &params, &res)
	_, err := ctx.RunResult()
	if err != nil {
		return "", err
	}
	return res.Hash, nil
}

func (client *JSONClient) GetBlockByHash(blockHash string) (*rpctypes.BlockOverview, error) {
	params := rpctypes.QueryParm{
		Hash: blockHash,
	}
	var res rpctypes.BlockOverview
	ctx := jsonclient.NewRPCCtx(client.url, "Chain33.GetBlockOverview", params, &res)
	_, err := ctx.RunResult()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (client *JSONClient) GetBalance(addrs []string, execer string) ([]*rpctypes.Account, error) {
	params := types.ReqBalance{
		Addresses: addrs,
		Execer:    execer,
		StateHash: "",
	}
	var res []*rpctypes.Account
	ctx := jsonclient.NewRPCCtx(client.url, "Chain33.GetBalance", params, &res)
	_, err := ctx.RunResult()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ValidateAddress(addr string) bool {
	err := address.CheckAddress(addr, -1)
	if err != nil {
		return false
	}
	return true
}
