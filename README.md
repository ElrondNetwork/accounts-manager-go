# elrond-accounts-manager-go

The go implementation for the elrond-accounts-manager

- This application will be responsible to fetch all elrond-accounts that have staked an amount of EGLD tokens. 
Alter the accounts are fetched from API will process all the information and will index the new data 
in a new Elaticsearch index.

- The new Elastisearch index will contain all the accounts that have balance and also information 
about the staked balance.

### Sources of accounts with stake

- This go client will fetch information from:
    1. Validators system smart contract
    2. Delegation manager system smart contracts
    3. Legacy delegation smart contract
    

### Installation and running


#### Step 1: install & configure go:

The installation of go should proceed as shown in official golang 
installation guide https://golang.org/doc/install . In order to run the node, minimum golang 
version should be 1.12.4.


#### Step 2: clone the repository and build the binary:

```
 $ git clone https://github.com/ElrondNetwork/accounts-manager-go.git
 $ cd accounts-manager-go/cmd/manager
 $ GO111MODULE=on go mod vendor
 $ go build
```

#### Step 3: run manager
```
 $ ./manager --config="pathToConfig/config.toml"
```
