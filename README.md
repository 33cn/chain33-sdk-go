# chain33-sdk-go
chain33 sdk golang

## 版本
golang1.13 or latest

## 安装

```text
//开启mod功能
export GO111MODULE=on

//国内用户需要导入阿里云代理，用于下载依赖包
export GOPROXY=https://mirrors.aliyun.com/goproxy
```

## 使用

### RPC客户端
通过jclient调用rpc接口发送交易和查询交易
```go
client, err := client.NewJSONClient(name, server)
client.SendTransaction(signedTx) 
```

### 加解密
crypto包实现了常用的加解密算法，签名算法，sha256哈希算法，密码生成和区块链地址生成，可通过crypto包直接调用。
```go
// 国密签名
sig, _ := crypto.gm.SM2Sign(priv, msg,nil)

// 国密验签
result := crypto.gm.SM2Verify(pub, msg, nil, sig)
```

### 代理重加密
代理重加密密钥生成和本地加解密

### 存证接口
创建存证合约原始交易

### EVM合约
- 合约部署
```go
code, err := types.FromHex(codes)
tx, err := CreateEvmContract(code, "", "evm-sdk-test", paraName)
err = SignTx(tx, deployPrivateKey)
signTx := types.ToHexPrefix(types.Encode(tx))
client.SendTransaction(signedTx)

contractAddress := crypto.GetExecAddress(deployAddress + strings.TrimPrefix(txhash, "0x")) 
```

- 合约调用
```go
param := fmt.Sprintf("mint(%s,%s,%s,%s)", useraAddress, idStr, amountStr, uriStr)
initNFT, err := EncodeParameter(abi, param)
tx, err = CallEvmContract(initNFT, "", 0, contractAddress, paraName)
// 构造代扣交易组
group, err := CreateNobalance(tx, useraPrivateKey, withholdPrivateKey, paraName)
signTx = types.ToHexPrefix(types.Encode(group.Tx()))
client.SendTransaction(signTx)
```

- 合约查询
```go
param = fmt.Sprintf("balanceOf(%s,%d)", useraAddress, ids[0])
QueryContract(url, contractAddress, abi, param, contractAddress)
```

## 接口文档
[chain33-sdk-go API](./接口文档.md)