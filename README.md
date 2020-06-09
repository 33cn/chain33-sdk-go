# chain33-sdk-go
chain33 sdk golang

### 接口文档

#### 1. 账户相关
#### 1.1 创建账户

**函数原型**
```
NewAccount(signType string) (*Account, error)
```

**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|signType|string|是|签名类型,支持"secp256k1",默认"secp256k1"|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|Account|PrivateKey  []byte | 私钥
|	    |PublicKey   []byte|  公钥
|       |Address     string | 地址
|       |SignType    string|  签名类型

#### 1.2 交易签名

**函数原型**
```
SignRawTransaction(raw string, privateKey string, signType string) (string, error)
```

**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|raw|string|是|原始交易数据|
|privateKey|string|是|私钥|
|signType|string|是|签名类型|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|result|string | 签名后的交易

#### 2. 代理重加密
#### 2.1 生成对称加密秘钥

**函数原型**
```
GenerateEncryptKey(pubOwner []byte) ([]byte, string, string)
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|pubOwner|[]byte|是|加密用户非对称秘钥公钥|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|key |[]byte | 对称加密秘钥
|pub_r |string | 随机公钥r，用于重加密授权
|pub_u |string | 随机公钥u，用于重加密授权

#### 2.2 生成重加密秘钥分片

**函数原型**
```
GenerateKeyFragments(privOwner []byte, pubRecipient []byte, numSplit, threshold int) ([]*KFrag, error) 
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|privOwner |[]byte|是|共享用户私钥|
|pubRecipient |[]byte|是|授权用户公钥|
|numSplit |int |是|秘钥分片数|
|threshold |int|是|最小秘钥分片重组阈值，不得大于numSplit|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|KFrag |Random    string  | 随机数，每个分片不同
| |Value     string | 重加密证明，每个分片不同
| |PrecurPub string | 随机公钥，所有分片相同
|err  |error | 错误返回

#### 2.3 重组重加密秘钥分片

**函数原型**
```
AssembleReencryptFragment(privRecipient []byte, reKeyFrags []*ReKeyFrag) ([]byte, error)
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|privRecipient |[]byte|是|授权用户私钥|
|reKeyFrags |ReKeyFrag|是|重加密秘钥分片，从各个重加密节点获取|
```
type ReKeyFrag struct {
	ReKeyR    string // 重加密证明R
	ReKeyU    string // 重加密证明U
	Random    string // 随机数，每个分片不同
	PrecurPub string // 随机公钥，所有分片相同
}
```

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|key |[]byte  | 重组后的对称秘钥
|err  |error | 错误返回

#### 3. RPC客户端
#### 3.1 创建jsonRPC客户端

**函数原型**
```
NewJSONClient(prefix, url string) (*JSONClient, error)
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|prefix |string|是|前缀|
|url |string|是|rpc服务端连接|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|JSONClient |JSONClient | 客户端对象
|err  |error | 错误返回

#### 3.2 rpc调用

**函数原型**
```
JSONClient.Call(method string, params, resp interface{}) error
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|method |string|是|rpc调用方法|
|params |interface|是|调用方法对应的参数|
|resp |interface|是|调用方法对应的返回值|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|err  |error | 错误返回

#### 3.3 发送交易

**函数原型**
```
JSONClient.SendTransaction(signedTx string) (string, error)
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|signedTx |string|是|已签名的交易|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|string  |string | 返回哈希
|err  |error | 错误返回

#### 3.4 交易查询

**函数原型**
```
JSONClient.QueryTransaction(hash string) (*TransactionDetail, error)
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|hash |string|是|交易哈希|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|TransactionDetail  |TransactionDetail | 交易详情
|err  |error | 错误返回
