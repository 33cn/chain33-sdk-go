# storage 合约sdk 使用

## 直接调用github.com/33cn/chain33-sdk-go/dapp/storage下面封装得函数

#### storage
#### 1.1 创建明文存证交易

**函数原型**
```
CreateContentStorageTx(paraName string, op int32, key string, content []byte, value string) (*types.Transaction, error)
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|paraName |string|是|平行链链名前缀，主链填空字符串|
|op |string|否|操作类型，0为新建存储，1为追加|
|key |string|否|不填默认采用txhash为key值|
|content |bytes|否|content和value二选一填值|
|value |string|否|content和value二选一填值|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|*types.Transaction|*types.Transaction | 未签名的交易
|err  |error | 错误返回

#### 1.2 创建hash存证交易

**函数原型**
```
CreateHashStorageTx(paraName string, key string, hash []byte, value string) (*types.Transaction, error)
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|paraName |string|是|平行链链名前缀，主链填空字符串|
|key |string|否|不填默认采用txhash为key值|
|hash |bytes|否|hash和value二选一填值|
|value |string|否|hash和value二选一填值|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|*types.Transaction|*types.Transaction | 未签名的交易
|err  |error | 错误返回

#### 1.3 创建链接存证交易

**函数原型**
```
CreateLinkStorageTx(paraName string, key string, link []byte, value string) (*types.Transaction, error) 
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|paraName |string|是|平行链链名前缀，主链填空字符串|
|key |string|否|不填默认采用txhash为key值|
|link |bytes|否|link和value二选一填值|
|value |string|否|link和value二选一填值|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|*types.Transaction|*types.Transaction | 未签名的交易
|err  |error | 错误返回

#### 1.4 创建隐私存证交易

**函数原型**
```
CreateEncryptStorageTx(paraName string, key string, contentHash, encryptContent, nonce []byte, value string) *types.Transaction
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|paraName |string|是|平行链链名前缀，主链填空字符串|
|key |string|否|不填默认采用txhash为key值|
|contentHash |bytes|是|明文hash值用于校验结果一致性|
|encryptContent |bytes|是|源文件的密文，由加密key及nonce对明文加密得到该值|
|nonce |bytes|是|加密iv，通过AES进行加密时制定随机生成的iv,解密时需要使用该值|
|value |string|否|预留字段|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|*types.Transaction|*types.Transaction | 未签名的交易
|err  |error | 错误返回

#### 1.4 创建分享隐私交易

**函数原型**
```
CreateEncryptShareStorageTx(paraName string, key string, contentHash, encryptContent, publickey []byte, value string) *types.Transaction
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|paraName |string|是|平行链链名前缀，主链填空字符串|
|key |string|否|不填默认采用txhash为key值|
|contentHash |bytes|是|明文hash值用于校验结果一致性|
|encryptContent |bytes|是|源文件得密文,采用非对称加密，用公钥地址加密，用私钥解密|
|publickey |bytes|是|公钥|
|value |string|否|预留字段|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|*types.Transaction|*types.Transaction | 未签名的交易
|err  |error | 错误返回

#### 1.5 根据key查询存储内容

**函数原型**
```
QueryStorageByKey(prefix, url, key string) (*types.Storage, error)
```
**请求参数**

|参数|类型|是否必填|说明|
|----|----|----|----|
|prefix |string|是|平行链链名前缀，主链填空字符串|
|url |string|是|地址|
|key |string|是|key值|

**返回字段：**

|返回字段|字段类型|说明|
|----|----|----|
|*types.Storage|*types.Storage| 存证内容
|err  |error | 错误返回

字段信息，可参考github.com/33cn/chain33-sdk-go/proto/sotrage.proto文件