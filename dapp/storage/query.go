package storage

import (
	"encoding/json"
	"fmt"
	. "github.com/bitly/go-simplejson"
	"github.com/33cn/chain33-sdk-go/client"
	"github.com/33cn/chain33-sdk-go/types"
)

func QueryStorageByKey(prefix, url, key string) (*types.Storage, error) {
	jsonClient, err := client.NewJSONClient(prefix, url)
	if err != nil {
		return nil, err
	}
	jsonraw, err := json.Marshal(&types.QueryStorage{TxHash: key})
	if err != nil {
		return nil, err
	}
	query := client.Query4Jrpc{
		Execer:   prefix + StorageX,
		FuncName: FuncNameQueryStorage,
		Payload:  jsonraw,
	}
	//var storage types.Storage
	data, err := jsonClient.CallBack("Chain33.Query", query, ParseStorage)
	if err != nil {
		return nil, err
	}
	return data.(*types.Storage), nil
}

//回调解析函数
func ParseStorage(raw json.RawMessage) (interface{}, error) {
	js, err := NewJson(raw)
	if err != nil {
		return nil, err
	}
	if contentStorge, ok := js.CheckGet("contentStorage"); ok {
		contentHex, _ := contentStorge.Get("content").String()
		content, _ := types.FromHex(contentHex)
		key, _ := contentStorge.Get("key").String()
		value, _ := contentStorge.Get("value").String()
		op, _ := contentStorge.Get("op").Int()
		storage := &types.Storage{Ty: TyContentStorageAction, Value: &types.Storage_ContentStorage{ContentStorage: &types.ContentOnlyNotaryStorage{
			Content:              content,
			Key:                  key,
			Op:                   int32(op),
			Value:                value,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		}}}
		return storage, nil
	}
	if linkStorge, ok := js.CheckGet("linkStorage"); ok {
		linkHex, _ := linkStorge.Get("link").String()
		link, _ := types.FromHex(linkHex)
		key, _ := linkStorge.Get("key").String()
		value, _ := linkStorge.Get("value").String()
		storage := &types.Storage{Ty: TyLinkStorageAction, Value: &types.Storage_LinkStorage{LinkStorage: &types.LinkNotaryStorage{
			Link:                 link,
			Key:                  key,
			Value:                value,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		}}}
		return storage, nil
	}
	if hashStorge, ok := js.CheckGet("hashStorage"); ok {
		hashHex, _ := hashStorge.Get("hash").String()
		hash, _ := types.FromHex(hashHex)
		key, _ := hashStorge.Get("key").String()
		value, _ := hashStorge.Get("value").String()
		storage := &types.Storage{Ty: TyHashStorageAction, Value: &types.Storage_HashStorage{HashStorage: &types.HashOnlyNotaryStorage{
			Hash:                 hash,
			Key:                  key,
			Value:                value,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		}}}
		return storage, nil
	}

	if encryptStorage, ok := js.CheckGet("encryptStorage"); ok {
		contentHashHex, _ := encryptStorage.Get("contentHash").String()
		contentHash, _ := types.FromHex(contentHashHex)
		encryptContentHex, _ := encryptStorage.Get("encryptContent").String()
		encryptContent, _ := types.FromHex(encryptContentHex)
		nonceHex, _ := encryptStorage.Get("nonce").String()
		nonce, _ := types.FromHex(nonceHex)
		key, _ := encryptStorage.Get("key").String()
		value, _ := encryptStorage.Get("value").String()
		storage := &types.Storage{Ty: TyEncryptStorageAction, Value: &types.Storage_EncryptStorage{EncryptStorage: &types.EncryptNotaryStorage{
			EncryptContent:       encryptContent,
			ContentHash:          contentHash,
			Nonce:                nonce,
			Key:                  key,
			Value:                value,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		}}}
		return storage, nil
	}
	if encryptStorage, ok := js.CheckGet("encryptShareStorage"); ok {
		contentHashHex, _ := encryptStorage.Get("contentHash").String()
		contentHash, _ := types.FromHex(contentHashHex)
		encryptContentHex, _ := encryptStorage.Get("encryptContent").String()
		encryptContent, _ := types.FromHex(encryptContentHex)
		pubKeyHex, _ := encryptStorage.Get("pubKey").String()
		pubKey, _ := types.FromHex(pubKeyHex)
		key, _ := encryptStorage.Get("key").String()
		value, _ := encryptStorage.Get("value").String()
		storage := &types.Storage{Ty: TyEncryptShareStorageAction, Value: &types.Storage_EncryptShareStorage{EncryptShareStorage: &types.EncryptShareNotaryStorage{
			EncryptContent:       encryptContent,
			ContentHash:          contentHash,
			PubKey:               pubKey,
			Key:                  key,
			Value:                value,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		}}}
		return storage, nil
	}
	return nil, fmt.Errorf("unknow type!")
}
