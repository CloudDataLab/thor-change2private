# Private Thor
rebuilding public blockchain project [Vechain/thor](https://github.com/vechain/thor) to a private chain, which based on the version [19 Jul, 2018 Vechain/thor:a9b9ea5](https://github.com/vechain/thor/tree/a9b9ea581bfa656bfd9a9a89fa7d8440ad7cf12c)



## what we change

- P2P accessing manage

We make thor to be a private chain mainly through P2P accessing as we maintain a white-list in builting smart contract.

- restricted transaction

We also restrict some users to sending transaction through builting smart contract. 

- blockInterval

The constant parameter **blockInterval ** controls the rate of block generation. We want to set it dynamically with smart contract.

## how to use

see the **branch test** and refer original **[Vechain/thor](https://github.com/vechain/thor)** project. Plz go to wiki to know more.