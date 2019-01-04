// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package genesis

import (
	"math/big"

	"github.com/vechain/thor/builtin"
	"github.com/vechain/thor/state"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
	"github.com/vechain/thor/vm"

)

// NewTestnet create genesis for testnet.
func NewTestnet() *Genesis {
	launchTime := uint64(1530014400) // 'Tue Jun 26 2018 20:00:00 GMT+0800 (CST)'

	// use this address as executor instead of builtin one, for test purpose
	executor, _ := thor.ParseAddress("0xcaE0d5a9D4380234bc562f9C409B418688dBC57E")
	acccount0, _ := thor.ParseAddress("0x9C9AF4fD05A876D2bF1f8A1dd22522cb57DF1987")

	master0, _ := thor.ParseAddress("0x532f30873af3bf0a40a1ea2de7785a1ddae2c82e")
	endorser0, _ := thor.ParseAddress("0xe59eD2484F2F63a8E99f9348ECa482b371A8a3F4")

	test0,_ := thor.ParseAddress("0x76c5117e049E78F53b73B5d7575e5E88A42f1447")

	builder := new(Builder).
		Timestamp(launchTime).
		GasLimit(thor.InitialGasLimit).
		State(func(state *state.State) error {
		tokenSupply := new(big.Int)

		// alloc precompiled contracts
		for addr := range vm.PrecompiledContractsByzantium {
			state.SetCode(thor.Address(addr), emptyRuntimeBytecode)
		}

		// setup builtin contracts
		state.SetCode(builtin.Authority.Address, builtin.Authority.RuntimeBytecodes())
		//edit by sion
		state.SetCode(builtin.Trader.Address,builtin.Trader.RuntimeBytecodes())
		state.SetCode(builtin.Energy.Address, builtin.Energy.RuntimeBytecodes())
		state.SetCode(builtin.Params.Address, builtin.Params.RuntimeBytecodes())
		state.SetCode(builtin.Prototype.Address, builtin.Prototype.RuntimeBytecodes())
		state.SetCode(builtin.Extension.Address, builtin.Extension.RuntimeBytecodes())

		// 50 billion for account0
		amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(50*1000*1000*1000))
		state.SetBalance(acccount0, amount)
		state.SetEnergy(acccount0, &big.Int{}, launchTime)
		tokenSupply.Add(tokenSupply, amount)

		// 25 million for endorser0
		amount = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(25*1000*1000))
		state.SetBalance(endorser0, amount)
		state.SetEnergy(endorser0, &big.Int{}, launchTime)
		tokenSupply.Add(tokenSupply, amount)

		amount = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(30*1000*1000))
		state.SetBalance(executor, amount)
		state.SetEnergy(executor, &big.Int{}, launchTime)
		tokenSupply.Add(tokenSupply, amount)

		amount = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(30*1000*1000))
		state.SetBalance(test0, amount)
		state.SetEnergy(test0, &big.Int{}, launchTime)
		tokenSupply.Add(tokenSupply, amount)

		builtin.Energy.Native(state, launchTime).SetInitialSupply(tokenSupply, &big.Int{})
		return nil
	}).
	// set initial params
	// use an external account as executor to manage testnet easily
		Call(
		tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", thor.KeyExecutorAddress, new(big.Int).SetBytes(executor[:]))),
		thor.Address{}).
		Call(
		tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", thor.KeyRewardRatio, thor.InitialRewardRatio)),
		executor).
		Call(
		tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", thor.KeyBaseGasPrice, thor.InitialBaseGasPrice)),
		executor).
		Call(
		tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", thor.KeyProposerEndorsement, thor.InitialProposerEndorsement)),
		executor).
		// add master0 as the initial block proposer
		Call(tx.NewClause(&builtin.Authority.Address).WithData(mustEncodeInput(builtin.Authority.ABI, "add", master0, endorser0, thor.BytesToBytes32([]byte("master0")),"112.74.47.241/32,120.78.83.87/32")),
		executor).
		Call(tx.NewClause(&builtin.Trader.Address).WithData(mustEncodeInput(builtin.Trader.ABI, "add", executor , thor.BytesToBytes32([]byte("master0")))),
		executor).
		Call(tx.NewClause(&builtin.Trader.Address).WithData(mustEncodeInput(builtin.Trader.ABI, "add", master0 , thor.BytesToBytes32([]byte("master0")))),
		executor).
		// by kasper
		// initial block-interval
		Call(
		tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", thor.KeyBlockInterval, thor.InitialBlockInterval)),
		executor)
	id, err := builder.ComputeID()
	if err != nil {
		panic(err)
	}
	return &Genesis{builder, id, "testnet"}
}
