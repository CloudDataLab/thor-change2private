// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package trader

import (
	"github.com/vechain/thor/thor"
)

type (
	entry struct {
		Identity thor.Bytes32
		Prev     *thor.Address `rlp:"nil"`
		Next     *thor.Address `rlp:"nil"`
	}

	// Candidate candidate of block proposer.
	TraderInfo struct {
		AccountTrader thor.Address
		Identity   thor.Bytes32
	}
)

// IsEmpty returns whether the entry can be treated as empty.
func (e *entry) IsEmpty() bool {
	return e.Identity.IsZero() &&
		e.Prev == nil &&
		e.Next == nil
}

func (e *entry) IsLinked() bool {
	return e.Prev != nil || e.Next != nil
}
