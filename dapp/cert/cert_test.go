package cert

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
	url = "http://127.0.0.1:8801"
	keyFilePath = "./test/keystore/5c3682a5719cf5bc1bd6280938670c3acfcb67cc15744a7b9b348066795a4e62_sk"
	certFilePath = "./test/signcerts/user1@org1-cert.pem"
)

func TestCreateCertNormalTx(t *testing.T) {
	account,err := sdk.NewAccountFromLocal(crypto.SM2, keyFilePath)
	assert.Nil(t, err)

	certByte,err := types.ReadFile(certFilePath)
	assert.Nil(t, err)

	tx, err := CreateCertNormalTx("", account.PrivateKey, certByte, []byte("cert test"), "key1", []byte("value1"))
	assert.Nil(t, err)

	jsonclient, err := client.NewJSONClient("", url)
	assert.Nil(t, err)

	signTx := types.ToHexPrefix(types.Encode(tx))
	reply, err := jsonclient.SendTransaction(signTx)
	assert.Nil(t, err)

	txhash := types.ToHexPrefix(sdk.Hash(tx))
	assert.Equal(t, txhash, reply)

	time.Sleep(2 * time.Second)
	detail, err := jsonclient.QueryTransaction(txhash)
	assert.Nil(t, err)
	assert.Equal(t, types.ExecOk, int(detail.Receipt.Ty))

}
