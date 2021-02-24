package storage

import (
	sdk "github.com/33cn/chain33-sdk-go"
	"github.com/33cn/chain33-sdk-go/client"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	privkey = "cc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"
	url     = "http://127.0.0.1:8801"
)

func TestCreateContentStorageTx(t *testing.T) {
	//第一次存储
	tx, err := CreateContentStorageTx("", OpCreate, "", []byte("hello"), "")
	assert.Nil(t, err)
	hexbytes, _ := types.FromHex(privkey)
	sdk.Sign(tx, hexbytes, crypto.SECP256K1, nil)
	txhash := types.ToHexPrefix(sdk.Hash(tx))
	jsonclient, err := client.NewJSONClient("", url)
	assert.Nil(t, err)
	signTx := types.ToHexPrefix(types.Encode(tx))
	reply, err := jsonclient.SendTransaction(signTx)
	assert.Nil(t, err)
	assert.Equal(t, txhash, reply)
	time.Sleep(2 * time.Second)
	detail, err := jsonclient.QueryTransaction(txhash)
	assert.Nil(t, err)
	assert.Equal(t, types.ExecOk, int(detail.Receipt.Ty))
	//查询
	storage, err := QueryStorageByKey("", url, txhash)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), storage.GetContentStorage().Content)
	//第二次追加  老版本不支持
	//tx,err=CreateContentStorageTx("",OpAdd,txhash,[]byte("world"),"")
	//assert.Nil(t,err)
	//tx.Sign(types.SECP256K1, priv)
	//signTx =common.ToHexPrefix(types.Encode(tx))
	//_,err=jsonclient.SendTransaction(signTx)
	//assert.Nil(t,err)
	//time.Sleep(2*time.Second)
	////查询
	//storage,err=QueryStorageByKey("",url,txhash)
	//assert.Nil(t,err)
	//assert.Equal(t,[]byte("hello,world"),storage.GetContentStorage().Content)

}

//hash,or link 存证
func TestCreateHashStorageTx(t *testing.T) {
	tx, err := CreateHashStorageTx("", "", []byte("123456harrylee"), "")
	assert.Nil(t, err)
	//签名
	hexbytes, _ := types.FromHex(privkey)
	sdk.Sign(tx, hexbytes, crypto.SECP256K1, nil)
	txhash := types.ToHexPrefix(sdk.Hash(tx))
	jsonclient, err := client.NewJSONClient("", url)
	assert.Nil(t, err)
	signTx := types.ToHexPrefix(types.Encode(tx))
	reply, err := jsonclient.SendTransaction(signTx)
	assert.Nil(t, err)
	assert.Equal(t, txhash, reply)
	time.Sleep(time.Second)
	detail, err := jsonclient.QueryTransaction(txhash)
	assert.Nil(t, err)
	assert.Equal(t, types.ExecOk, int(detail.Receipt.Ty))
	//查询
	storage, err := QueryStorageByKey("", url, txhash)
	assert.Nil(t, err)
	assert.Equal(t, []byte("123456harrylee"), storage.GetHashStorage().Hash)
}

//hash,or link 存证
func TestCreateLinkStorageTx(t *testing.T) {
	tx, err := CreateLinkStorageTx("", "", []byte("hello"), "")
	assert.Nil(t, err)
	hexbytes, _ := types.FromHex(privkey)
	sdk.Sign(tx, hexbytes, crypto.SECP256K1, nil)
	txhash := types.ToHexPrefix(sdk.Hash(tx))
	jsonclient, err := client.NewJSONClient("", url)
	assert.Nil(t, err)
	signTx := types.ToHexPrefix(types.Encode(tx))
	reply, err := jsonclient.SendTransaction(signTx)
	assert.Nil(t, err)
	assert.Equal(t, txhash, reply)
	time.Sleep(time.Second)
	storage, err := QueryStorageByKey("", url, txhash)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), storage.GetLinkStorage().Link)
}

func TestByteFromHex(t *testing.T) {
	hex := "0x313233343536"
	data, err := types.FromHex(hex)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))

}
