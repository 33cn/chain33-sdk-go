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
	caPrivKey=""
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
	caAdminPriv, _ := types.FromHex(caPrivKey)

	userName1 := "ycy"
	identity1 := "101"
	pubKey1   := "pubKey1"

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

// 证书管理方案测试
func TestCAServerAdmin(t *testing.T) {
    // CA中心管理员
	caAdminPriv, _ := types.FromHex(caPrivKey)

	// 测试用户
	userName2 := "ycy2"
	identity2 := "102"
	pubKey2   := "pubKey2"

	// 测试证书管理员
	userName3 := "ycy3"
	privKey3,_  := types.FromHex("privKey3")
	pubKey3   := "pubKey3"

	// 测试非证书管理员
	userName4  := "ycy4"
	privKey4,_ := types.FromHex("privKey4")

	jsonclient, err := client.NewJSONClient("", caUrl)
	assert.Nil(t, err)

	// 注册用户
	res1, err := jsonclient.CertUserRegister(userName2, identity2, pubKey2, "", caAdminPriv)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, true, res1)
	fmt.Printf("注册用户：%t\n", res1)

	// 非管理员申请证书，返回失败
	res2, err := jsonclient.CertEnroll(identity2, userName4, privKey4)
	assert.NotNil(t, err)
	assert.Nil(t, res2)
	fmt.Printf("普通用户申请证书：%s， 错误：%s\n", res2, err.Error())

	// 添加管理员
	res3, err := jsonclient.CertAdminRegister(userName3, pubKey3, caAdminPriv)
	assert.Equal(t, true, res3)
	fmt.Printf("添加管理员：%t\n", res3)

	// 新管理员申请证书, 返回成功
	res4, err := jsonclient.CertEnroll(identity2, userName3, privKey3)
	assert.Nil(t, err)
	assert.NotNil(t, res4)
	fmt.Printf("管理员申请证书：%s\n", res4.Serial)

	// 非管理员注销证书，返回失败
	res5, err := jsonclient.CertRevoke(res4.Serial, "", userName4, privKey4)
	assert.NotNil(t, err)
	assert.Equal(t, false, res5)
	fmt.Printf("普通用户注销证书：%t， 错误：%s\n", res5, err.Error())

	// 新管理员注销证书, 返回成功
	res6, err := jsonclient.CertRevoke(res4.Serial, "", userName3, privKey3)
	assert.Nil(t, err)
	assert.Equal(t, true, res6)
	fmt.Printf("管理员注销证书：%t\n", res6)

}

// 批量生成多个用户证书
func TestCAServerGenerate(t *testing.T) {
	caAdminPriv, _ := types.FromHex(caPrivKey)

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

// 生成新证书，修改userName和identity
func TestCAServerGenerateNew(t *testing.T) {
	caAdminPriv, _ := types.FromHex(caPrivKey)

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

// 注销新证书，修改待注销的证书对应用户的identity
func TestCAServerRevokeNew(t *testing.T) {
	caAdminPriv, _ := types.FromHex(caPrivKey)

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