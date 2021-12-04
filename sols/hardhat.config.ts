import { task } from "hardhat/config";
import "@nomiclabs/hardhat-waffle";
import "@nomiclabs/hardhat-ethers";
import "@nomiclabs/hardhat-etherscan";
import "@openzeppelin/hardhat-upgrades";
import { HardhatUserConfig } from "hardhat/config";

task("accounts", "Prints the list of accounts", async (_, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(`${account.address} -> ${await account.getBalance()}`);
  }
});

const pkey = process.env.ACCOUNT_PRIVATE_KEYS?.split(",") || [];

const net = {
  ropsten: {
    url: `${process.env.ROPSTEN_ALCHEMYAPI}`,
    accounts: pkey,
  },
  kovan: {
    url: `${process.env.KOVAN_ALCHEMYAPI}`,
    accounts: pkey,
  },
  rinkeby: {
    url: `${process.env.RINKEBY_ALCHEMYAPI}`,
    accounts: pkey,
  },
};

const conf: HardhatUserConfig = {
  solidity: "0.8.10",
  networks: net,
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY,
  },
};

export default conf;
