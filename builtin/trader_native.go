// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/xenv"
)

func init() {
	defines := []struct {
		name string
		run  func(env *xenv.Environment) []interface{}
	}{
		{"native_executor", func(env *xenv.Environment) []interface{} {
			env.UseGas(thor.SloadGas)
			addr := thor.BytesToAddress(Params.Native(env.State()).Get(thor.KeyExecutorAddress).Bytes())
			return []interface{}{addr}
		}},
		{"native_add", func(env *xenv.Environment) []interface{} {
			var args struct {
				AccountTrader common.Address
				Identity   common.Hash
			}
			env.ParseArgs(&args)

			env.UseGas(thor.SloadGas)
			ok := Trader.Native(env.State()).Add(
				thor.Address(args.AccountTrader),
				thor.Bytes32(args.Identity))

			if ok {
				env.UseGas(thor.SstoreSetGas)
				env.UseGas(thor.SstoreResetGas)
			}
			return []interface{}{ok}
		}},
		{"native_revoke", func(env *xenv.Environment) []interface{} {
			var accountTrader common.Address
			env.ParseArgs(&accountTrader)

			env.UseGas(thor.SloadGas)
			ok := Trader.Native(env.State()).Revoke(thor.Address(accountTrader))
			if ok {
				env.UseGas(thor.SstoreResetGas * 3)
			}
			return []interface{}{ok}
		}},
		{"native_get", func(env *xenv.Environment) []interface{} {
			var accountTrader common.Address
			env.ParseArgs(&accountTrader)

			env.UseGas(thor.SloadGas * 2)

			listed, identity := Trader.Native(env.State()).Get(thor.Address(accountTrader))

			return []interface{}{listed, identity}
		}},
		{"native_first", func(env *xenv.Environment) []interface{} {
			env.UseGas(thor.SloadGas)
			if accountTrader := Trader.Native(env.State()).First(); accountTrader != nil {
				return []interface{}{*accountTrader}
			}
			return []interface{}{thor.Address{}}
		}},
		{"native_next", func(env *xenv.Environment) []interface{} {
			var accountTrader common.Address
			env.ParseArgs(&accountTrader)

			env.UseGas(thor.SloadGas)
			if next := Trader.Native(env.State()).Next(thor.Address(accountTrader)); next != nil {
				return []interface{}{*next}
			}
			return []interface{}{thor.Address{}}
		}},
	}
	abi := Trader.NativeABI()
	for _, def := range defines {
		if method, found := abi.MethodByName(def.name); found {
			nativeMethods[methodKey{Trader.Address, method.ID()}] = &nativeMethod{
				abi: method,
				run: def.run,
			}
		} else {
			panic("method not found: " + def.name)
		}
	}
}
