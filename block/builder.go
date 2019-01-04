// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package block

import (
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
)

// Builder to make it easy to build a block object.
type Builder struct {
	headerBody headerBody
	txs        tx.Transactions
}

// ParentID set parent id.
func (b *Builder) ParentID(id thor.Bytes32) *Builder {
	b.headerBody.ParentID = id
	return b
}

// Timestamp set timestamp.
func (b *Builder) Timestamp(ts uint64) *Builder {
	b.headerBody.Timestamp = ts
	return b
}

// TotalScore set total score.
func (b *Builder) TotalScore(score uint64) *Builder {
	b.headerBody.TotalScore = score
	return b
}

// GasLimit set gas limit.
func (b *Builder) GasLimit(limit uint64) *Builder {
	b.headerBody.GasLimit = limit
	return b
}

// GasUsed set gas used.
func (b *Builder) GasUsed(used uint64) *Builder {
	b.headerBody.GasUsed = used
	return b
}

// Beneficiary set recipient of reward.
func (b *Builder) Beneficiary(addr thor.Address) *Builder {
	b.headerBody.Beneficiary = addr
	return b
}

// StateRoot set state root.
func (b *Builder) StateRoot(hash thor.Bytes32) *Builder {
	b.headerBody.StateRoot = hash
	return b
}

// ReceiptsRoot set receipts root.
func (b *Builder) ReceiptsRoot(hash thor.Bytes32) *Builder {
	b.headerBody.ReceiptsRoot = hash
	return b
}

// Transaction add a transaction.
func (b *Builder) Transaction(tx *tx.Transaction) *Builder {
	b.txs = append(b.txs, tx)
	return b
}

// Build build a block object.
func (b *Builder) Build() *Block {
	header := Header{body: b.headerBody}
	header.body.TxsRoot = b.txs.RootHash()

	return &Block{
		header: &header,
		txs:    b.txs,
	}
}