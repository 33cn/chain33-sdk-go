package eth

import (
	"testing"

	"github.com/33cn/chain33-sdk-go/types"
	"github.com/33cn/chain33/common"
	"github.com/stretchr/testify/assert"
)

func TestEth(t *testing.T) {
	privStr := "3967abcafaea83fee72766ca6dae578f4f156b5d1dae1ddf119e4564d5e2658c"
	addr := "0x6856f610b40e7321cace9e1f8752315110862573"

	priv, err := types.FromHex(privStr)
	assert.Nil(t, err)
	pub := PubKeyFromPrivate(priv)
	calAddr, err := PubKeyToAddress(pub)
	assert.Nil(t, err)
	assert.Equal(t, addr, calAddr)
	msg := common.Sha256([]byte("test eth"))
	sig := Sign(msg, priv)
	res := Validate(msg, pub, sig)
	assert.Equal(t, true, res)
}
