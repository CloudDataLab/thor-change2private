package node

import (
	"bytes"
	"math/rand"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/vechain/thor/genesis"
	"github.com/vechain/thor/lvldb"
	"github.com/vechain/thor/tx"
)

func newTx() *tx.Transaction {
	tx := new(tx.Builder).Nonce(rand.Uint64()).Build()
	sig, _ := crypto.Sign(tx.SigningHash().Bytes(), genesis.DevAccounts()[0].PrivateKey)
	return tx.WithSignature(sig)
}

func TestTxStash(t *testing.T) {
	db, _ := lvldb.NewMem()
	defer db.Close()

	stash := newTxStash(db, 10)

	var saved tx.Transactions
	for i := 0; i < 11; i++ {
		tx := newTx()
		assert.Nil(t, stash.Save(tx))
		saved = append(saved, tx)
	}

	loaded := newTxStash(db, 10).LoadAll()

	saved = saved[1:]
	sort.Slice(saved, func(i, j int) bool {
		return bytes.Compare(saved[i].ID().Bytes(), saved[j].ID().Bytes()) < 0
	})

	sort.Slice(loaded, func(i, j int) bool {
		return bytes.Compare(loaded[i].ID().Bytes(), loaded[j].ID().Bytes()) < 0
	})

	assert.Equal(t, saved.RootHash(), loaded.RootHash())
}
