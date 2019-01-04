'use strict';

import { thorify} from "thorify";
import { Params } from "./params";

const Web3 = require("web3");		// recommand use require() instead of import here
const web3 = new Web3();
thorify(web3, "http://localhost:8669");

//add private key to wallet
//private key of executor
web3.eth.accounts.wallet.add("0xc97fd1e5296fb9d53bb179f093ad458194f12f8419f8bda86c8602f9807b90ce")
//private key of ....f1447(public key)
web3.eth.accounts.wallet.add("0x226f7869edb01a713521f993aa83132de567f6d3a83d7007fc58d47dd1557305")
//private key of ....2c82e
web3.eth.accounts.wallet.add("0xa8c8a78979003a3a0bdc163b0228fd75b771b9c19db6b2f9aa2761a6d76fc381")
//private key of ....f1987
web3.eth.accounts.wallet.add("0x10180c258862575cf35709bd3210a1314c1fa5411730f0144f183e9cf0dcd5f2")

const params = new web3.eth.Contract(Params.contractABI,Params.contractAddress)

params.methods.get("0x000000000000000000000000000000000000626c6f636b2d696e74657276616c").call("latest",function(err:any,result:any){
  console.log(result)
})
