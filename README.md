# ETH Bridge Zone

[![CircleCI](https://circleci.com/gh/cosmos/peggy/tree/master.svg?style=svg)](https://circleci.com/gh/cosmos/peggy/tree/master)

## Summary

Unidirectional Peggy is the starting point for cross chain value transfers from the Ethereum blockchain to Cosmos-SDK based blockchains as part of the Ethereum Cosmos Bridge project. The system accepts incoming transfers of Ethereum tokens on an Ethereum smart contract, locking them while the transaction is validated and equitable funds issued to the intended recipient on the Cosmos bridge chain.

## Disclaimer

This codebase, including all smart contract components, have not been professionally audited and are not intended for use in a production environment. As such, users should NOT trust the system to securely hold mainnet funds. Any developers attempting to use Unidirectional Peggy on the mainnet at this time will need to develop their own smart contracts or find another implementation.

## Architecture

See [here](./docs/architecture.md)

## Requirements
 - Go 1.13

## Example application

These modules can be added to any Cosmos-SDK based chain, but a demo application/blockchain is provided with example code for how to integrate them. It can be installed and built as follows:

```
# Clone the repository
mkdir -p $GOPATH/src/github.com/cosmos
cd $GOPATH/src/github.com/cosmos
git clone https://github.com/cosmos/peggy
cd peggy && git checkout master

# Install tools (golangci-lint v1.18)
make tools-clean
make tools

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands, confirming the build is successful:
ebd help
ebcli help
ebrelayer help
```

## Running and testing the application

First, initialize a chain and create accounts to test sending of a random token.

```bash
# Initialize the genesis.json file that will help you to bootstrap the network
ebd init local --chain-id=peggy

# Create a key to hold your validator account and for another test account
ebcli keys add validator
# Enter password

ebcli keys add testuser
# Enter password

# Initialize the genesis account and transaction
ebd add-genesis-account $(ebcli keys show validator -a) 1000000000stake,1000000000atom

# Create genesis transaction
ebd gentx --name validator
# Enter password

# Collect genesis transaction
ebd collect-gentxs

# Now its safe to start `ebd`
ebd start

# Then, wait 10 seconds and in another terminal window, test things are ok by sending 10 tok tokens from the validator to the testuser
ebcli tx send validator $(ebcli keys show testuser -a) 10stake --chain-id=peggy --yes

# Wait a few seconds for confirmation, then confirm token balances have changed appropriately
ebcli query account $(ebcli keys show validator -a) --trust-node
ebcli query account $(ebcli keys show testuser -a) --trust-node

# See the help for the ethbridge create claim function
ebcli tx ethbridge create-claim --help

# Now you can test out the ethbridge module by submitting a claim for an ethereum prophecy
# Create a bridge claim (Ethereum prophecies are stored on the blockchain with an identifier created by concatenating the nonce and sender address)
ebcli tx ethbridge create-claim 0 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 3eth --from=validator --chain-id=peggy --yes

# Then read the prophecy to confirm it was created with the claim added
ebcli query ethbridge prophecy 0 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 --trust-node

# And finally, confirm that the prophecy was successfully processed and that new eth was minted to the testuser address
ebcli query account $(ebcli keys show testuser -a) --trust-node

```

## Using the application from rest-server

First, run the cli rest-server

```bash
ebcli rest-server --trust-node
```

An api collection for Postman (https://www.getpostman.com/) is provided [here](./docs/peggy.postman_collection.json) which documents some API endpoints and can be used to interact with it.
Note: For checking account details/balance, you will need to change the cosmos addresses in the URLs, params and body to match the addresses you generated that you want to check.

## Running the relayer service

For automated relaying, there is a relayer service that can be run that will automatically watch and relay events.

```bash
# Check ebrelayer connection to ebd
ebrelayer status

# Initialize the Relayer service for automatic claim processing
ebrelayer init wss://ropsten.infura.io/ws ec6df30846baab06fce9b1721608853193913c19 "LogLock\(bytes32,address,bytes,address,uint256,uint256\)" validator --chain-id=peggy

# Enter password and press enter
# You should see a message like:  Started ethereum websocket... and Subscribed to contract events...
```

The relayer will now watch the contract on Ropsten and create a claim whenever it detects a lock event.

## Using the bridge

With the application set up and the relayer running, you can now use Peggy by sending a lock transaction to the smart contract. You can do this from any Ethereum wallet/client that supports smart contract transactions.

### Set up

Create a .env file with variables MNEMONIC, INFURA_PROJECT_ID and LOCAL_PROVIDER. An example configuration can be found in .env.example. For running the bridge locally, you'll only need the LOCAL_PROVIDER. For running the bridge on ropsten testnet, you'll need the MNEMONIC from MetaMask and the INFURA_PROJECT_ID from Infura.

### Terminal 1: Start local blockchain

```
$ cd testnet-contracts/
$ yarn develop
```

### Terminal 2: Compile, deploy, check contract's deployed address

```
$ cd testnet-contracts/

# Copy contract ABI to go modules:
$ yarn peggy:abi

# Deploy contract to local blockchain
$ yarn migrate

# Get contract's address
$ yarn peggy:address
```

### Terminal 3: Build and start Ethereum Bridge

```
# Build the Ethereum Bridge application
$ make install

# Start the Bridge's blockchain
$ ebd start
```

### Terminal 4: Start the relayer service

Example [LOCAL_WEB_SOCKET]: ws://127.0.0.1:8545/
Example [PEGGY_DEPLOYED_ADDRESS]: 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB

```
# Start ebrelayer on the contract's deployed address
$ ebrelayer init [LOCAL_WEB_SOCKET] [PEGGY_DEPLOYED_ADDRESS] LogLock\(bytes32,address,bytes,address,uint256,uint256\) validator --chain-id=testing

# Enter password and press enter
# You should see a message like: Started ethereum websocket with provider: [LOCAL_WEB_SOCKET] \ Subscribed to contract events on address: [PEGGY_DEPLOYED_ADDRESS]
```

### Using Terminal 2: Send lock transaction to contract

The lock transaction uses the default parameters:

- [HASHED_COSMOS_RECIPIENT_ADDRESS] = 0x636f736d6f7331706a74677530766175326d35326e72796b64707a74727438383761796b756530687137646668
- [DEPLOYED_TOKEN_ADDRESS] = 0x0000000000000000000000000000000000000000
- [WEI_AMOUNT] = 10

```
$ yarn peggy:lock

# Expected successful output in the relayer console:

New Lock Transaction:
Tx hash: 0x83e6ee88c20178616e68fee2477d21e84f16dcf6bac892b18b52c000345864c0
Block number: 5
Event ID: cc10955295e555130c865949fb1fd48dba592d607ae582b43a2f3f0addce83f2
Token: 0x0000000000000000000000000000000000000000
Sender: 0xc230f38FF05860753840e0d7cbC66128ad308B67
Recipient: cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh
Value: 10
Nonce: 1

# You can also confirm the tokens have been minted by using the CLI again:
$ ebcli query account cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh --trust-node
```

### Running on the testnet

To run the Ethereum Bridge on the testnet, repeat the steps for running locally except for the following changes:

```
# Specify the ropsten network via a --network flag for the following commands...
$ yarn migrate --network ropsten
$ yarn peggy:address --network ropsten
$ yarn peggy:lock --network ropsten

# Start ebrelayer with ropsten network websocket
$ ebrelayer init wss://ropsten.infura.io/ [PEGGY_DEPLOYED_ADDRESS] LogLock\(bytes32,address,bytes,address,uint256,uint256\) validator --chain-id=testing
```

## Using the modules in other projects

The ethbridge and oracle modules can be used in other cosmos-sdk applications by copying them into your application's modules folders and including them in the same way as in the example application. Each module may be moved to its own repo or integrated into the core Cosmos-SDK in future, for easier usage.

For instructions on building and deploying the smart contracts, see the README in their folder.
