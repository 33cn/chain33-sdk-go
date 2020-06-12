package crypto

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)


func Sha256(b []byte) []byte {
	hasher := sha256.New()
	hasher.Write(b)
	return hasher.Sum(nil)
}

func Sha2Sum(b []byte) []byte {
	tmp := sha256.Sum256(b)
	tmp = sha256.Sum256(tmp[:])
	return tmp[:]
}

func rimpHash(in []byte, out []byte) {
	sha := sha256.New()
	_, err := sha.Write(in)
	if err != nil {
		return
	}
	rim := ripemd160.New()
	_, err = rim.Write(sha.Sum(nil)[:])
	if err != nil {
		return
	}
	copy(out, rim.Sum(nil))
}

// Rimp160 Returns hash: RIMP160( SHA256( data ) )
// Where possible, using RimpHash() should be a bit faster
func Rimp160(b []byte) []byte {
	out := make([]byte, 20)
	rimpHash(b, out[:])
	return out[:]
}

func KDF(x []byte, length int) []byte {
	var c []byte

	var ct int64 = 1
	h := sha256.New()
	for i, j := 0, (length+31)/32; i < j; i++ {
		h.Reset()
		h.Write(x)
		h.Write(big.NewInt(ct).Bytes())
		hash := h.Sum(nil)
		if i+1 == j && length%32 != 0 {
			c = append(c, hash[:length%32]...)
		} else {
			c = append(c, hash...)
		}
		ct++
	}

	return c
}