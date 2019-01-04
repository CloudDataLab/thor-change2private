// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

pragma solidity 0.4.24;

/// @title Authority manages a candidates list of master nodes(block proposers).
contract Trader {
    function executor() public view returns(address) {
        return TraderNative(this).native_executor();
    }

    function add(address _accountTrader, bytes32 _identity) public {
        require(_accountTrader != 0, "builtin: invalid account trader");
        require(_identity != 0, "builtin: invalid identity");
        require(msg.sender == executor(), "builtin: executor required");

        require(TraderNative(this).native_add(_accountTrader, _identity), "builtin: already exists");

        emit Manage(_accountTrader, "added");
    }

    function revoke(address _accountTrader) public {
        require(msg.sender == executor(), "builtin: requires executor");
        require(TraderNative(this).native_revoke(_accountTrader), "builtin: not listed");

        emit Manage(_accountTrader, "revoked");
    }

    function get(address _accountTrader) public view returns(bool listed, bytes32 identity) {
        return TraderNative(this).native_get(_accountTrader);
    }

    function first() public view returns(address) {
        return TraderNative(this).native_first();
    }

    function next(address _accountTrader) public view returns(address) {
        return TraderNative(this).native_next(_accountTrader);
    }

    event Manage(address indexed accountTrader, bytes32 action);
}

contract TraderNative {
    function native_executor() public view returns(address);
    function native_add(address accountTrader, bytes32 identity) public returns(bool);
    function native_revoke(address accountTrader) public returns(bool);
    function native_get(address accountTrader) public view returns(bool, bytes32);
    function native_first() public view returns(address);
    function native_next(address accountTrader) public view returns(address);
}
