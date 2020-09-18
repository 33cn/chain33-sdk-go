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

## 接口文档
[chain33-sdk-go API](./接口文档.md)