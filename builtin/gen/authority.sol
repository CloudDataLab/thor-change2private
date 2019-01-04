// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

pragma solidity 0.4.24;

/// @title Authority manages a candidates list of master nodes(block proposers).
contract Authority {
    function executor() public view returns(address) {
        return AuthorityNative(this).native_executor();
    }
    //edit by sion
    function add(address _nodeMaster, address _endorsor, bytes32 _identity,string _nodeIp) public {
        require(_nodeMaster != 0, "builtin: invalid node master");
        require(_endorsor != 0, "builtin: invalid endorsor");
        require(_identity != 0, "builtin: invalid identity");
        require(bytes(_nodeIp).length != 0, "builtin: invalid nodeIp");//edit by sion
        require(msg.sender == executor(), "builtin: executor required");

        require(AuthorityNative(this).native_add(_nodeMaster, _endorsor, _identity, _nodeIp), "builtin: already exists");//edit by sion

        emit Candidate(_nodeMaster, "added");
    }

    function revoke(address _nodeMaster) public {
        require(msg.sender == executor() || !AuthorityNative(this).native_isEndorsed(_nodeMaster), "builtin: requires executor, or node master out of endorsed");
        require(AuthorityNative(this).native_revoke(_nodeMaster), "builtin: not listed");

        emit Candidate(_nodeMaster, "revoked");
    }
    //edit by sion
    function get(address _nodeMaster) public view returns(bool listed, address endorsor, bytes32 identity, string nodeIp, bool active) {
        return AuthorityNative(this).native_get(_nodeMaster);
    }

    function first() public view returns(address) {
        return AuthorityNative(this).native_first();
    }

    function next(address _nodeMaster) public view returns(address) {
        return AuthorityNative(this).native_next(_nodeMaster);
    }

    event Candidate(address indexed nodeMaster, bytes32 action);
}

contract AuthorityNative {
    function native_executor() public view returns(address);
    //edit by sion
    function native_add(address nodeMaster, address endorsor, bytes32 identity, string nodeIp) public returns(bool);
    function native_revoke(address nodeMaster) public returns(bool);
    function native_get(address nodeMaster) public view returns(bool, address, bytes32, string, bool);//edit by sion
    function native_first() public view returns(address);
    function native_next(address nodeMaster) public view returns(address);
    function native_isEndorsed(address nodeMaster) public view returns(bool);
}
