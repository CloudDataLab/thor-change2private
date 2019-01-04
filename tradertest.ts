'use strict';

import { thorify} from "thorify";
import { Trader } from "./trader";

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

//the object of contract trader
const trader = new web3.eth.Contract(Trader.contractABI,"0x0000000000000000000000000000547261646572")

// console.log(web3.utils.asciiToHex("master1").)

// trader.methods.first().call("latest",function(err:any,result:any){
//   console.log(result)
// })

// trader.methods.next("0xe59eD2484F2F63a8E99f9348ECa482b371A8a3F4").call("latest",function(err:any,result:any){
//   console.log(result)
// })

//add a new trader account
// trader.methods.add("0x76c5117e049E78F53b73B5d7575e5E88A42f1447","0x000000000000000000000000000000000000000000000000006d617374657231").send({
//   from : "0xcaE0d5a9D4380234bc562f9C409B418688dBC57E",gas:1090000
// })

//revoke a trader account
// trader.methods.revoke("0xe59eD2484F2F63a8E99f9348ECa482b371A8a3F4").send({
// from : "0xcaE0d5a9D4380234bc562f9C409B418688dBC57E",gas : 1090000
// })

//get information by trader's account
// trader.methods.get("0x76c5117e049E78F53b73B5d7575e5E88A42f1447").call("latest",function(err:any,result:any){
//   console.log(result)
// })



//send transaction
// web3.eth.sendTransaction({
//    from: "0x76c5117e049E78F53b73B5d7575e5E88A42f1447",
//    to: "0xcaE0d5a9D4380234bc562f9C409B418688dBC57E",
//    value: 10,
//  })

//get balance by address of account
// web3.eth.getBalance("0x76c5117e049E78F53b73B5d7575e5E88A42f1447").then(function(result : any){
//     console.log("balance is: ", result)
// })
