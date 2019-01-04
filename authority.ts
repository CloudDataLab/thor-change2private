// origin abstracts the smart contract related contents:
export class Authority
{
    static contractAddress = "0x0000000000000000000000417574686f72697479";
    static contractABI = [
      {
        "anonymous": false,
        "inputs": [
          {
            "indexed": true,
            "name": "nodeMaster",
            "type": "address"
          },
          {
            "indexed": false,
            "name": "action",
            "type": "bytes32"
          }
        ],
        "name": "Candidate",
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
            "name": "_nodeMaster",
            "type": "address"
          },
          {
            "name": "_endorsor",
            "type": "address"
          },
          {
            "name": "_identity",
            "type": "bytes32"
          },
          {
            "name": "_nodeIp",
            "type": "string"
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
            "name": "_nodeMaster",
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
            "name": "_nodeMaster",
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
            "name": "endorsor",
            "type": "address"
          },
          {
            "name": "identity",
            "type": "bytes32"
          },
          {
            "name": "nodeIp",
            "type": "string"
          },
          {
            "name": "active",
            "type": "bool"
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
            "name": "_nodeMaster",
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
