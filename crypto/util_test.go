package crypto

import (
	"testing"

	ttypes "github.com/33cn/chain33/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc1, err := NewAccount(BTC_ADDRESS)
	assert.Nil(t, err)
	tx1 := &ttypes.Transaction{Execer: []byte("none"), Payload: []byte("btc test")}
	err = SignTx(tx1, acc1.PrivateKey, acc1.Type)
	assert.Nil(t, err)
	ret1 := tx1.CheckSign(-1)
	assert.True(t, ret1)

	acc2, err := NewAccount(ETH_ADDRESS)
	assert.Nil(t, err)
	tx2 := &ttypes.Transaction{Execer: []byte("none"), Payload: []byte("eth test")}
	err = SignTx(tx2, acc2.PrivateKey, acc2.Type)
	assert.Nil(t, err)
	ret2 := tx2.CheckSign(-1)
	assert.True(t, ret2)
}
