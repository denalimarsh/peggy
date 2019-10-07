require("dotenv").config();

var HDWalletProvider = require("@truffle/hdwallet-provider");

module.exports = {
  networks: {
    develop: {
      host: "localhost",
      port: 7545, // Match default network 'ganache'
      network_id: "*",
      gas: 6721975, // Truffle default development block gas limit
      gasPrice: 200000000000,
      solc: {
        version: "0.5.0",
        optimizer: {
          enabled: true,
          runs: 200
        }
      }
    },
    ropsten: {
      provider: function() {
        return new HDWalletProvider(
          process.env.MNEMONIC,
          "https://ropsten.infura.io/".concat(process.env.INFURA_PROJECT_ID)
        );
      },
      network_id: 3,
      gas: 6000000
    }
  },
  coverage: {
    host: "localhost",
    network_id: "*",
    port: 8555,
    gas: 0xfffffffffff, // <-- Use this high gas value
    gasPrice: 0x01, // <-- Use this low gas price
    copyPackages: ["openzeppelin-solidity"]
  },
  rpc: {
    host: "localhost",
    post: 8080
  },
  mocha: {
    useColors: true
  }
};
