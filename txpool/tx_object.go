// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package txpool

import (
	"math/big"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/vechain/thor/block"
	"github.com/vechain/thor/chain"
	"github.com/vechain/thor/runtime"
	"github.com/vechain/thor/state"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
	"github.com/vechain/thor/builtin"
)

type txObject struct {
	*tx.Transaction
	resolved *runtime.ResolvedTransaction

	timeAdded       int64
	executable      bool
	overallGasPrice *big.Int // don't touch this value, it's only be used in pool's housekeeping
}

func resolveTx(tx *tx.Transaction) (*txObject, error) {
	resolved, err := runtime.ResolveTransaction(tx)
	if err != nil {
		return nil, err
	}

	return &txObject{
		Transaction: tx,
		resolved:    resolved,
		timeAdded:   time.Now().UnixNano(),
	}, nil
}

func (o *txObject) Origin() thor.Address {
	return o.resolved.Origin
}

func (o *txObject) Executable(chain *chain.Chain, state *state.State, headBlock *block.Header) (bool, error) {
	// by kasper
	// get block Interval
	blockInterval := builtin.Params.Native(state).Get(thor.KeyBlockInterval).Uint64()

	switch {
	case o.Gas() > headBlock.GasLimit():
		return false, errors.New("gas too large")
	case o.IsExpired(headBlock.Number()):
		return false, errors.New("expired")
	// by kasper
	//case o.BlockRef().Number() > headBlock.Number()+uint32(3600*24/thor.BlockInterval):
	case o.BlockRef().Number() > headBlock.Number()+uint32(3600*24/blockInterval):
		return false, errors.New("block ref out of schedule")
	}

	//edit by sion
	var trader=builtin.Trader.Native(state)
	signer,_:=o.Signer()
	if listed,_:=trader.Get(signer);!listed{
		return false,errors.New("builtin: invalid account")
	}
	//edit by sion

	if _, err := chain.GetTransactionMeta(o.ID(), headBlock.ID()); err != nil {
		if !chain.IsNotFound(err) {
			return false, err
		}
	} else {
		return false, errors.New("known tx")
	}

	if dep := o.DependsOn(); dep != nil {
		txMeta, err := chain.GetTransactionMeta(*dep, headBlock.ID())
		if err != nil {
			if chain.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		if txMeta.Reverted {
			return false, errors.New("dep reverted")
		}
	}

	if o.BlockRef().Number() > headBlock.Number() {
		return false, nil
	}

	checkpoint := state.NewCheckpoint()
	defer state.RevertTo(checkpoint)

	// by kasper
	//if _, _, _, _, err := o.resolved.BuyGas(state, headBlock.Timestamp()+thor.BlockInterval); err != nil {
	if _, _, _, _, err := o.resolved.BuyGas(state, headBlock.Timestamp()+blockInterval); err != nil {
		return false, err
	}
	return true, nil
}

func sortTxObjsByOverallGasPriceDesc(txObjs []*txObject) {
	sort.Slice(txObjs, func(i, j int) bool {
		gp1, gp2 := txObjs[i].overallGasPrice, txObjs[j].overallGasPrice
		return gp1.Cmp(gp2) >= 0
	})
}

