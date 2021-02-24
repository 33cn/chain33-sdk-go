package crypto

import "github.com/33cn/chain33-sdk-go/types"

func EncodeCertToSignature(signBytes, cert, uid []byte) []byte {
	var certSignature types.CertSignature
	certSignature.Cert = cert
	certSignature.Signature = signBytes
	certSignature.Uid = uid

	return types.Encode(&certSignature)
}
