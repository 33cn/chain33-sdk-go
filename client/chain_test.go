package client

import (
	"fmt"
	"testing"

	sdk "github.com/33cn/chain33-sdk-go"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	url := "http://127.0.0.1:8801"
	client, err := NewJSONClient("", url)
	assert.Nil(t, err)

	header, err := client.GetLastHeader()
	assert.Nil(t, err)
	fmt.Println("当前最大区块高度为:", header.Height)

	blocks, err := client.GetBlocks(1, 3, true)
	assert.Nil(t, err)
	for i := 0; i < len(blocks.Items); i++ {
		block := blocks.Items[i].Block
		fmt.Println("前一个区块hash:", block.ParentHash)
		fmt.Println("默克尔根hash:", block.TxHash)
		fmt.Println("区块高度:", block.Height)
		fmt.Println("区块时间戳:", block.BlockTime)
		fmt.Println("区块中交易数:", len(block.Txs))
		for j := 0; j < len(block.Txs); j++ {
			fmt.Printf("第%d笔交易详情:%v\n", j+1, block.Txs[j])
		}
	}

	blockHash, err := client.GetBlockHashByHeight(2)
	assert.Nil(t, err)
	fmt.Println("指定高度的区块hash值是:", blockHash)

	blockview, err := client.GetBlockByHash(blockHash)
	assert.Nil(t, err)
	fmt.Println("区块hash:", blockview.Head.Hash)
	fmt.Println("前一个区块hash:", blockview.Head.ParentHash)
	fmt.Println("默克尔根hash:", blockview.Head.TxHash)
	fmt.Println("区块高度:", blockview.Head.Height)
	fmt.Println("区块时间戳:", blockview.Head.BlockTime)
	fmt.Println("区块中交易数:", blockview.Head.TxCount)

	addrs := []string{"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt", "17RH6oiMbUjat3AAyQeifNiACPFefvz3Au"}
	accounts, err := client.GetBalance(addrs, "coins")
	for _, acc := range accounts {
		fmt.Println(acc.Addr, "balance is:", acc.Balance/1e8)
	}
}

func TestUtil(t *testing.T) {
	account, err := sdk.NewAccount("")
	assert.Nil(t, err)
	fmt.Println("new address:", account.Address)
	fmt.Println("new private key:", types.ToHex(account.PrivateKey))

	address := "1G1L2M1w1c1gpV6SP8tk8gBPGsJe2RfTks"
	isValid := ValidateAddress(address)
	fmt.Println("validate result is:", isValid)
}
