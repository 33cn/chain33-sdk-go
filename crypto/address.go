package crypto

import (
	"crypto/sha256"
	"fmt"
	"github.com/mr-tron/base58/base58"
)

//不同币种的前缀版本号
var coinPrefix = map[string][]byte{
	"BTC":  {0x00},
	"BCH":  {0x00},
	"BTY":  {0x00},
	"LTC":  {0x30},
	"ZEC":  {0x1c, 0xb8},
	"USDT": {0x00},
}

var addrSeed = []byte("address seed bytes for public key")

//MaxExecNameLength 执行器名最大长度
const MaxExecNameLength = 100

func PubKeyToAddress(pub []byte) (addr string, err error) {
	if len(pub) != 33 && len(pub) != 65 { //压缩格式 与 非压缩格式
		return "", fmt.Errorf("invalid public key byte")
	}

	//添加版本号
	hash160res := append(coinPrefix["BTY"], Rimp160(pub)...)

	//添加校验码
	cksum := checksum(hash160res)
	address := append(hash160res, cksum[:]...)

	//地址进行base58编码
	addr = base58.Encode(address)
	return
}

//checksum: first four bytes of double-SHA256.
func checksum(input []byte) (cksum [4]byte) {
	h := sha256.New()
	_, err := h.Write(input)
	if err != nil {
		return
	}
	intermediateHash := h.Sum(nil)
	h.Reset()
	_, err = h.Write(intermediateHash)
	if err != nil {
		return
	}
	finalHash := h.Sum(nil)
	copy(cksum[:], finalHash[:])
	return
}

func GetExecAddress(name string) string {
	if len(name) > MaxExecNameLength {
		panic("name too long")
	}
	var bname [200]byte
	buf := append(bname[:0], addrSeed...)
	buf = append(buf, []byte(name)...)
	pub := Sha2Sum(buf)

	var ad [25]byte
	ad[0] = 0
	copy(ad[1:21], Rimp160(pub))

	sh := Sha2Sum(ad[0:21])
	checksum := make([]byte, 4)
	copy(checksum, sh[:4])

	copy(ad[21:25], checksum[:])
	addr := base58.Encode(ad[:])

	return addr
}