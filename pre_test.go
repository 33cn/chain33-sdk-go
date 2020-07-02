package sdk

import (
	"fmt"
	"github.com/33cn/chain33-sdk-go/client"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ReqSendKeyFragment struct {
	PubOwner             string   `protobuf:"bytes,1,opt,name=pubOwner,proto3" json:"pubOwner,omitempty"`
	PubRecipient         string   `protobuf:"bytes,2,opt,name=pubRecipient,proto3" json:"pubRecipient,omitempty"`
	PubProofR            string   `protobuf:"bytes,3,opt,name=pubProofR,proto3" json:"pubProofR,omitempty"`
	PubProofU            string   `protobuf:"bytes,4,opt,name=pubProofU,proto3" json:"pubProofU,omitempty"`
	Random               string   `protobuf:"bytes,5,opt,name=random,proto3" json:"random,omitempty"`
	Value                string   `protobuf:"bytes,6,opt,name=value,proto3" json:"value,omitempty"`
	Expire               int64    `protobuf:"varint,7,opt,name=expire,proto3" json:"expire,omitempty"`
	DhProof              string   `protobuf:"bytes,8,opt,name=dhProof,proto3" json:"dhProof,omitempty"`
	PrecurPub            string   `protobuf:"bytes,9,opt,name=precurPub,proto3" json:"precurPub,omitempty"`
}

type ReqReeencryptParam struct {
	PubOwner             string   `protobuf:"bytes,1,opt,name=pubOwner,proto3" json:"pubOwner,omitempty"`
	PubRecipient         string   `protobuf:"bytes,2,opt,name=pubRecipient,proto3" json:"pubRecipient,omitempty"`
}

type RepReeencrypt struct {
	ReKeyR               string   `protobuf:"bytes,1,opt,name=reKeyR,proto3" json:"reKeyR,omitempty"`
	ReKeyU               string   `protobuf:"bytes,2,opt,name=reKeyU,proto3" json:"reKeyU,omitempty"`
	Random               string   `protobuf:"bytes,3,opt,name=random,proto3" json:"random,omitempty"`
	PrecurPub            string   `protobuf:"bytes,4,opt,name=precurPub,proto3" json:"precurPub,omitempty"`
}

func TestPre(t *testing.T) {
	privOwner,_ := types.FromHex("6d52c4680c00dcdb9d904dc6878a8e1c753ecf9c43a48499d819fdc0eafa4639")
	pubOwner,_ := types.FromHex("02e5fdf78aded517e3235c2276ed0e020226c55835dea7b8306f2e8d3d99d2d4f4")
	serverPub,_ := types.FromHex("02005d3a38feaff00f1b83014b2602d7b5b39506ddee7919dd66539b5428358f08")
	privRecipient, _ := types.FromHex("841e3b4ab211eecfccb475940171150fd1536cb656c870fe95d206ebf9732b6c")
	pubRecipient, _ := types.FromHex("03b9d801f88c38522a9bf786f23544259d516ee0d1f6699f926f891ac3fb92c6d9")
	msg := "hello proxy-re-encrypt"
	serverList := []string {"http://192.168.0.155:11801", "http://192.168.0.155:11802", "http://192.168.0.155:11803"}

	enKey, pub_r, pub_u := GeneratePreEncryptKey(pubOwner)
	cipher, err := crypto.AESCBCPKCS7Encrypt(enKey, []byte(msg))
	if err != nil {
		panic(err)
	}
	fmt.Println(types.ToHex(cipher))

	if err != nil {
		panic(err)
	}
	keyFrags, err := GenerateKeyFragments(privOwner, pubRecipient, 3, 2)
	if err != nil {
		panic(err)
	}

	dhproof := types.ECDH(crypto.PrivateFromByte(privOwner), crypto.PublicFromByte(serverPub))
	for i, server := range serverList {
		jclient, err := client.NewJSONClient("Pre", server)
		if err != nil {
			panic(err)
		}

		var result interface{}
		param := &ReqSendKeyFragment{
			PubOwner:     types.ToHex(pubOwner),
			PubRecipient: types.ToHex(pubRecipient),
			PubProofR:    pub_r,
			PubProofU:    pub_u,
			Random:       keyFrags[i].Random,
			Value:        keyFrags[i].Value,
			Expire:       1000000,
			DhProof:      types.ToHex(dhproof),
			PrecurPub:    keyFrags[i].PrecurPub,
		}
		jclient.Call("CollectFragment", param, &result)
	}

	param := &ReqReeencryptParam{
		PubOwner:             types.ToHex(pubOwner),
		PubRecipient:         types.ToHex(pubRecipient),
	}

	var rekeys = make([]*ReKeyFrag, 2)
	for i:=0; i < 2; i++ {
		jclient, err := client.NewJSONClient("Pre", serverList[i])
		if err != nil {
			panic(err)
		}
		var result RepReeencrypt
		jclient.Call("Reencrypt", param, &result)

		rekeys[i] = &ReKeyFrag{
			ReKeyR:    result.ReKeyR,
			ReKeyU:    result.ReKeyU,
			Random:    result.Random,
			PrecurPub: result.PrecurPub,
		}
	}
	encKey,err := AssembleReencryptFragment(privRecipient, rekeys)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, enKey, encKey)

	msgDecrypt, err := crypto.AESCBCPKCS7Decrypt(encKey, cipher)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msgDecrypt))
}
