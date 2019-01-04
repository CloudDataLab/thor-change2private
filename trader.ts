// origin abstracts the smart contract related contents:
export class Trader
{
    static contractABI =   [
      {
        "anonymous": false,
        "inputs": [
          {
            "indexed": true,
            "name": "accountTrader",
            "type": "address"
          },
          {
            "indexed": false,
            "name": "action",
            "type": "bytes32"
          }
        ],
        "name": "Manage",
        "type": "event"
      },
      {
        "constant": true,
        "inputs": [],
        "name": "executor",
        "outputs": [
          {
            "name": "",
            "type": "address"
          }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
      },
      {
        "constant": false,
        "inputs": [
          {
            "name": "_accountTrader",
            "type": "address"
          },
          {
            "name": "_identity",
            "type": "bytes32"
          }
        ],
        "name": "add",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "constant": false,
        "inputs": [
          {
            "name": "_accountTrader",
            "type": "address"
          }
        ],
        "name": "revoke",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "constant": true,
        "inputs": [
          {
            "name": "_accountTrader",
            "type": "address"
          }
        ],
        "name": "get",
        "outputs": [
          {
            "name": "listed",
            "type": "bool"
          },
          {
            "name": "identity",
            "type": "bytes32"
          }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
      },
      {
        "constant": true,
        "inputs": [],
        "name": "first",
        "outputs": [
          {
            "name": "",
            "type": "address"
          }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
      },
      {
        "constant": true,
        "inputs": [
          {
            "name": "_accountTrader",
            "type": "address"
          }
        ],
        "name": "next",
        "outputs": [
          {
            "name": "",
            "type": "address"
          }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
      }
    ]
}
