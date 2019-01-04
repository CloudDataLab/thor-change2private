// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package trader

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vechain/thor/state"
	"github.com/vechain/thor/thor"
)

var (
	headKey = thor.Blake2b([]byte("head"))
	tailKey = thor.Blake2b([]byte("tail"))
)

// Authority implements native methods of `Authority` contract.
type Trader struct {
	addr  thor.Address
	state *state.State
}

// New create a new instance.
func New(addr thor.Address, state *state.State) *Trader {
	return &Trader{addr, state}
}

func (a *Trader) getEntry(accountTrader thor.Address) *entry {
	var entry entry
	a.state.DecodeStorage(a.addr, thor.BytesToBytes32(accountTrader[:]), func(raw []byte) error {
		if len(raw) == 0 {
			return nil
		}
		return rlp.DecodeBytes(raw, &entry)
	})
	return &entry
}

func (a *Trader) setEntry(accountTrader thor.Address, entry *entry) {
	a.state.EncodeStorage(a.addr, thor.BytesToBytes32(accountTrader[:]), func() ([]byte, error) {
		if entry.IsEmpty() {
			return nil, nil
		}
		return rlp.EncodeToBytes(entry)
	})
}

func (a *Trader) getAddressPtr(key thor.Bytes32) (addr *thor.Address) {
	a.state.DecodeStorage(a.addr, key, func(raw []byte) error {
		if len(raw) == 0 {
			return nil
		}
		return rlp.DecodeBytes(raw, &addr)
	})
	return
}

func (a *Trader) setAddressPtr(key thor.Bytes32, addr *thor.Address) {
	a.state.EncodeStorage(a.addr, key, func() ([]byte, error) {
		if addr == nil {
			return nil, nil
		}
		return rlp.EncodeToBytes(addr)
	})
}

// Get get candidate by node master address.
//edit by sion
func (a *Trader) Get(accountTrader thor.Address) (listed bool, identity thor.Bytes32) {
	entry := a.getEntry(accountTrader)
	if entry.IsLinked() {
		return true,  entry.Identity
	}
	// if it's the only node, IsLinked will be false.
	// check whether it's the head.
	ptr := a.getAddressPtr(headKey)
	listed = ptr != nil && *ptr == accountTrader
	return listed, entry.Identity
}

// Add add a new candidate.
//edit by sion
func (a *Trader) Add(accountTrader thor.Address, identity thor.Bytes32) bool {
	entry := a.getEntry(accountTrader)
	if !entry.IsEmpty() {
		return false
	}

	entry.Identity = identity
	tailPtr := a.getAddressPtr(tailKey)
	entry.Prev = tailPtr

	a.setAddressPtr(tailKey, &accountTrader)
	if tailPtr == nil {
		a.setAddressPtr(headKey, &accountTrader)
	} else {
		tailEntry := a.getEntry(*tailPtr)
		tailEntry.Next = &accountTrader
		a.setEntry(*tailPtr, tailEntry)
	}

	a.setEntry(accountTrader, entry)
	return true
}

// Revoke revoke candidate by given node master address.
// The entry is not removed, but set unlisted and inactive.
func (a *Trader) Revoke(accountTrader thor.Address) bool {
	entry := a.getEntry(accountTrader)
	if !entry.IsLinked() {
		return false
	}

	if entry.Prev == nil {
		a.setAddressPtr(headKey, entry.Next)
	} else {
		prevEntry := a.getEntry(*entry.Prev)
		prevEntry.Next = entry.Next
		a.setEntry(*entry.Prev, prevEntry)
	}

	if entry.Next == nil {
		a.setAddressPtr(tailKey, entry.Prev)
	} else {
		nextEntry := a.getEntry(*entry.Next)
		nextEntry.Prev = entry.Prev
		a.setEntry(*entry.Next, nextEntry)
	}

	entry.Next = nil
	entry.Prev = nil     // unlist
	entry.Identity = [32]byte{}
	a.setEntry(accountTrader, entry)
	return true
}


// Candidates picks a batch of candidates up to limit, that satisfy given endorsement.
func (a *Trader) Traders() []*TraderInfo {
	ptr := a.getAddressPtr(headKey)
	traders := make([]*TraderInfo, 0)
	for ptr != nil {
		entry := a.getEntry(*ptr)
		traders = append(traders, &TraderInfo{
			AccountTrader: *ptr,
			Identity:   entry.Identity,
		})

		ptr = entry.Next
	}
	return traders
}


// First returns node master address of first entry.
func (a *Trader) First() *thor.Address {
	return a.getAddressPtr(headKey)
}

// Next returns address of next node master address after given node master address.
func (a *Trader) Next(accountTrader thor.Address) *thor.Address {
	return a.getEntry(accountTrader).Next
}
