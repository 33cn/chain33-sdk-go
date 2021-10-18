package client

import (
	"github.com/33cn/chain33-sdk-go/crypto/secp256r1"
	"github.com/33cn/chain33-sdk-go/types"
)

// 注册用户
func (client *JSONClient) CertUserRegister(userName, identity, userPub, admin string, adminKey []byte) (bool, error) {
	send := &types.ReqRegisterUser{
		Name:                 userName,
		Identity:             identity,
		PubKey:               userPub,
		Admin:                admin,
	}
	msg := types.Encode(send)
	sign, err := secp256r1.Sign(msg, adminKey)
	if err != nil {
		return false, err
	}
	send.Sign = sign

	var res bool
	err = client.Call("chain33-ca-server.RegisterUser", send, &res)
	if err != nil {
		return false, err
	}

	return res, nil
}

// 注销用户
func (client *JSONClient) CertUserRevoke(identity, admin string, adminKey []byte) (bool, error) {
	send := &types.ReqRevokeUser{
		Identity:             identity,
		Admin:                admin,
	}

	msg := types.Encode(send)
	sign, err := secp256r1.Sign(msg, adminKey)
	if err != nil {
		return false, err
	}
	send.Sign = sign

	var res bool
	err = client.Call("chain33-ca-server.RevokeUser", send, &res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (client *JSONClient) CertEnroll(identity, admin string, adminKey []byte) (*types.RepEnroll, error) {
	send := &types.ReqEnroll{
		Identity:             identity,
		Admin:                admin,
	}

	msg := types.Encode(send)
	sign, err := secp256r1.Sign(msg, adminKey)
	if err != nil {
		return nil, err
	}
	send.Sign = sign

	var res types.RepEnroll
	err = client.Call("chain33-ca-server.Enroll", send, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}


func (client *JSONClient) CertRevoke(serial, identity, admin string, adminKey []byte) (bool, error) {
	send := &types.ReqRevokeCert{
		Serial:               serial,
		Identity:             identity,
		Admin:                admin,
	}

	msg := types.Encode(send)
	sign, err := secp256r1.Sign(msg, adminKey)
	if err != nil {
		return false, err
	}
	send.Sign = sign

	var res bool
	err = client.Call("chain33-ca-server.RevokeCert", send, &res)
	if err != nil {
		return false, err
	}

	return res, nil
}


func (client *JSONClient) CertGetCertInfo(serial string, adminKey []byte) (*types.RepGetCertInfo, error) {
	send := &types.ReqGetCertInfo{
		Sn:                   serial,
	}

	msg := types.Encode(send)
	sign, err := secp256r1.Sign(msg, adminKey)
	if err != nil {
		return nil, err
	}
	send.Sign = sign

	var res types.RepGetCertInfo
	err = client.Call("chain33-ca-server.GetCertInfo", send, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *JSONClient) CertGetUserInfo(identity string, adminKey []byte) (*types.RepGetCertInfo, error) {
	send := &types.ReqGetUserInfo{
		Identity:                   identity,
	}

	msg := types.Encode(send)
	sign, err := secp256r1.Sign(msg, adminKey)
	if err != nil {
		return nil, err
	}
	send.Sign = sign

	var res types.RepGetCertInfo
	err = client.Call("chain33-ca-server.GetUserInfo", send, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}


func (client *JSONClient) CertAdminRegister(userName, userPub string, adminKey []byte) (bool, error) {
	send := &types.ReqAdmin{
		Name:                 userName,
		PubKey:               userPub,
	}

	msg := types.Encode(send)
	sign, err := secp256r1.Sign(msg, adminKey)
	if err != nil {
		return false, err
	}
	send.Sign = sign

	var res bool
	err = client.Call("chain33-ca-server.AddCertAdmin", send, &res)
	if err != nil {
		return false, err
	}

	return res, nil
}


func (client *JSONClient) CertAdminRemove(userName, userPub string, adminKey []byte) (bool, error) {
	send := &types.ReqAdmin{
		Name:                 userName,
		PubKey:               userPub,
	}

	msg := types.Encode(send)
	sign, err := secp256r1.Sign(msg, adminKey)
	if err != nil {
		return false, err
	}
	send.Sign = sign

	var res bool
	err = client.Call("chain33-ca-server.RemoveCertAdmin", send, &res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (client *JSONClient) CertValidate(serials []string) ([]string, error) {
	send := &types.ReqValidateCert{
		Serials: serials,
	}
	var detail []string
	err := client.Call("chain33-ca-server.Validate", send, &detail)
	if err != nil {
		return nil, err
	}

	return detail, nil
}


