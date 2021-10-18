package cert

import (
	"fmt"
	sdk "github.com/33cn/chain33-sdk-go"
	"github.com/33cn/chain33-sdk-go/client"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/crypto/secp256r1"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

var (
	url = "http://192.168.203.128:8801"
	keyFilePath = "./test/keystore/5c3682a5719cf5bc1bd6280938670c3acfcb67cc15744a7b9b348066795a4e62_sk"
	certFilePath = "./test/signcerts/user1@org1-cert.pem"

	caUrl = "http://127.0.0.1:11901"
)

func TestCreateCertNormalTx(t *testing.T) {
	account,err := sdk.NewAccountFromLocal(crypto.SM2, keyFilePath)
	assert.Nil(t, err)

	//certByte,err := types.ReadFile(certFilePath)
	//assert.Nil(t, err)

	tx, err := CreateCertNormalTx("", account.PrivateKey, nil,nil, "key1", []byte("value1"))
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

func TestCAServer(t *testing.T) {
	caAdminPriv, _ := types.FromHex("e1d61ee8d20558b2c272589e9fe636c4c969e06f103c29dbf2b5a385f20a91e8")
	//caAdminPub := "02fc5356da98ce1f97c7bda404f162b765c30c497445ac28857fdbfab04c6f589c"

	userName1 := "ycy"
	identity1 := "101"
	pubKey1   := "02d5c86531a6bfa14d31afe986b49ea0cf0be395925304b757cede832874a11a3d"

	jsonclient, err := client.NewJSONClient("", caUrl)
	assert.Nil(t, err)

	// 注册用户
	res1, err := jsonclient.CertUserRegister(userName1, identity1, pubKey1, "", caAdminPriv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, true, res1)

	// 申请证书
	res2, err := jsonclient.CertEnroll(identity1, "", caAdminPriv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	fmt.Println(res2.Serial)
	_ = ioutil.WriteFile("user1.cert", res2.Cert, 666)
	_ = ioutil.WriteFile("user1.key", res2.Key, 666)

	// 证书校验,返回成功
	res3, err := jsonclient.CertValidate([]string{res2.Serial})
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, 0, len(res3))

	// 证书注销
	res4, err := jsonclient.CertRevoke(res2.Serial, "", "", caAdminPriv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, true, res4)

	// 证书校验，返回失败
	res5, err := jsonclient.CertValidate([]string{res2.Serial})
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, 1, len(res5))
}

func TestCAServerAdmin(t *testing.T) {
	caAdminPriv, _ := types.FromHex("e1d61ee8d20558b2c272589e9fe636c4c969e06f103c29dbf2b5a385f20a91e8")
	//caAdminPub := "02fc5356da98ce1f97c7bda404f162b765c30c497445ac28857fdbfab04c6f589c"

	userName2 := "ycy2"
	identity2 := "102"
	pubKey2   := "037dcc61f5bf3bbe67846e9f3ed50250c6a2ac33069ca07338dbb653034e9e3a7f"

	userName3 := "ycy3"
	privKey3,_  := types.FromHex("36c597b95a438ce2782db6d7a8812bf3fe5c85677c49c150946f091c1dc641ea")
	pubKey3   := "03e810079431c75c969ea8dded0e1e31db2ac5582a4484b9b22dd92f74ab9a9523"

	jsonclient, err := client.NewJSONClient("", caUrl)
	assert.Nil(t, err)

	// 注册用户
	res1, err := jsonclient.CertUserRegister(userName2, identity2, pubKey2, "", caAdminPriv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, true, res1)

	// 非管理员申请证书
	res2, err := jsonclient.CertEnroll(identity2, userName3, privKey3)
	assert.NotNil(t, err)
	assert.Nil(t, res2)

	// 添加管理员
	res3, err := jsonclient.CertAdminRegister(userName3, pubKey3, caAdminPriv)
	assert.Equal(t, true, res3)

	// 新管理员申请证书
	res4, err := jsonclient.CertEnroll(identity2, userName3, privKey3)
	assert.Nil(t, err)
	assert.NotNil(t, res4)
}

func TestCAServerGenerate(t *testing.T) {
	caAdminPriv, _ := types.FromHex("e1d61ee8d20558b2c272589e9fe636c4c969e06f103c29dbf2b5a385f20a91e8")
	//caAdminPub := "02fc5356da98ce1f97c7bda404f162b765c30c497445ac28857fdbfab04c6f589c"

	jsonclient, err := client.NewJSONClient("", caUrl)
	assert.Nil(t, err)

	certNum := 4
	for i := 0; i < certNum; i++ {
		userName := "user" + strconv.Itoa(i)
		identity := "10" + strconv.Itoa(i)
		priv, err := secp256r1.GeneratePrivateKey()
		if err != nil {
			assert.Fail(t, err.Error())
		}
		pub := secp256r1.PubKeyFromPrivate(priv)
		fmt.Println(userName)
		fmt.Println(identity)
		fmt.Println(types.ToHex(pub))
		// 注册用户
		res1, err := jsonclient.CertUserRegister(userName, identity, types.ToHex(pub), "", caAdminPriv)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		assert.Equal(t, true, res1)

		// 申请证书
		res2, err := jsonclient.CertEnroll(identity, "", caAdminPriv)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		fmt.Println(res2.Serial)
		_ = ioutil.WriteFile("user" + strconv.Itoa(i)+".cert", res2.Cert, 666)
		_ = ioutil.WriteFile("user" + strconv.Itoa(i)+".key", res2.Key, 666)
	}

}

func TestCAServerGenerateNew(t *testing.T) {
	caAdminPriv, _ := types.FromHex("e1d61ee8d20558b2c272589e9fe636c4c969e06f103c29dbf2b5a385f20a91e8")
	//caAdminPub := "02fc5356da98ce1f97c7bda404f162b765c30c497445ac28857fdbfab04c6f589c"

	jsonclient, err := client.NewJSONClient("", caUrl)
	assert.Nil(t, err)

	userName := "newUser"
	identity := "new101"
	priv, err := secp256r1.GeneratePrivateKey()
	if err != nil {
		assert.Fail(t, err.Error())
	}
	pub := secp256r1.PubKeyFromPrivate(priv)
	fmt.Println(userName)
	fmt.Println(identity)
	fmt.Println(types.ToHex(pub))

	// 注册用户
	res1, err := jsonclient.CertUserRegister(userName, identity, types.ToHex(pub), "", caAdminPriv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, true, res1)

	// 申请证书
	res2, err := jsonclient.CertEnroll(identity, "", caAdminPriv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	fmt.Println(res2.Serial)
	_ = ioutil.WriteFile("newUser.cert", res2.Cert, 666)
	_ = ioutil.WriteFile("newUser.key", res2.Key, 666)

}

func TestCAServerRevokeNew(t *testing.T) {
	caAdminPriv, _ := types.FromHex("e1d61ee8d20558b2c272589e9fe636c4c969e06f103c29dbf2b5a385f20a91e8")
	//caAdminPub := "02fc5356da98ce1f97c7bda404f162b765c30c497445ac28857fdbfab04c6f589c"

	jsonclient, err := client.NewJSONClient("", caUrl)
	assert.Nil(t, err)

	identity := "new101"

	// 证书注销
	res, err := jsonclient.CertRevoke("", identity, "", caAdminPriv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, true, res)

}