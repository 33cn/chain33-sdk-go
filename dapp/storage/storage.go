package storage

import (
	"fmt"
	"github.com/33cn/chain33-sdk-go/crypto"
	. "github.com/33cn/chain33-sdk-go/types"
	"math/rand"
)

/*
  //Op 0表示创建 1表示追加add (新版本才支持）
  int32  op = 1;
  //存在内容
  bytes content = 2;
  //自定义的主键，可以为空，如果没传，则用txhash为key
  string key = 3;
  //字符串值
  string value = 4;
  value 为预留字段
  content和value只能有一个有效值
*/
//明文存证溯源
func CreateContentStorageTx(paraName string, op int32, key string, content []byte, value string) (*Transaction, error) {
	if op != 0 && op != 1 {
		return nil, fmt.Errorf("unknow op..,op only 0 or 1,please check!")
	}
	payload := &StorageAction{Ty: TyContentStorageAction, Value: &StorageAction_ContentStorage{&ContentOnlyNotaryStorage{
		Content:              content,
		Op:                   op,
		Key:                  key,
		Value:                value,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}},
	}
	if paraName == "" {
		tx := &Transaction{Execer: []byte(StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: Addr}
		return tx, nil
	} else {
		tx := &Transaction{Execer: []byte(paraName + StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: crypto.GetExecAddress(paraName + StorageX)}
		return tx, nil
	}
}

//链接存证 paraName 平行链前缀，如果是主链则为空字符串，key唯一性，如果不填默认采用txhash为key
func CreateLinkStorageTx(paraName string, key string, link []byte, value string) (*Transaction, error) {
	payload := &StorageAction{Ty: TyLinkStorageAction, Value: &StorageAction_LinkStorage{&LinkNotaryStorage{
		Link:                 link,
		Key:                  key,
		Value:                value,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}},
	}
	if paraName == "" {
		tx := &Transaction{Execer: []byte(StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: Addr}
		return tx, nil
	} else {
		tx := &Transaction{Execer: []byte(paraName + StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: crypto.GetExecAddress(paraName + StorageX)}
		return tx, nil
	}
}

//hash存证
func CreateHashStorageTx(paraName string, key string, hash []byte, value string) (*Transaction, error) {
	payload := &StorageAction{Ty: TyHashStorageAction, Value: &StorageAction_HashStorage{&HashOnlyNotaryStorage{
		Hash:                 hash,
		Key:                  key,
		Value:                value,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}},
	}
	if paraName == "" {
		tx := &Transaction{Execer: []byte(StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: Addr}
		return tx, nil
	} else {
		tx := &Transaction{Execer: []byte(paraName + StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: crypto.GetExecAddress(paraName + StorageX)}
		return tx, nil
	}
}

//隐私加密存储
func CreateEncryptStorageTx(paraName string, key string, contentHash, encryptContent, nonce []byte, value string) *Transaction {
	payload := &StorageAction{Ty: TyEncryptStorageAction, Value: &StorageAction_EncryptStorage{&EncryptNotaryStorage{
		ContentHash:          contentHash,
		EncryptContent:       encryptContent,
		Nonce:                nonce,
		Key:                  key,
		Value:                value,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}},
	}
	if paraName == "" {
		tx := &Transaction{Execer: []byte(StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: Addr}
		return tx
	} else {
		tx := &Transaction{Execer: []byte(paraName + StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: crypto.GetExecAddress(paraName + StorageX)}
		return tx
	}
}

//分享隐私加密存储
func CreateEncryptShareStorageTx(paraName string, key string, contentHash, encryptContent, publickey []byte, value string) *Transaction {
	payload := &StorageAction{Ty: TyEncryptShareStorageAction, Value: &StorageAction_EncryptShareStorage{&EncryptShareNotaryStorage{
		ContentHash:          contentHash,
		EncryptContent:       encryptContent,
		PubKey:               publickey,
		Key:                  key,
		Value:                value,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}},
	}
	if paraName == "" {
		tx := &Transaction{Execer: []byte(StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: Addr}
		return tx
	} else {
		tx := &Transaction{Execer: []byte(paraName + StorageX), Payload: Encode(payload), Fee: 1e5, Nonce: rand.Int63(), To: crypto.GetExecAddress(paraName + StorageX)}
		return tx
	}
}
