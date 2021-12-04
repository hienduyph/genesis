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

const net = {
  ropsten: {
    url: `${process.env.ROPSTEN_ALCHEMYAPI}`,
    accounts: [`0x${process.env.ACCOUNT_PRIVATE_KEY}`],
  },
  kovan: {
    url: `${process.env.KOVAN_ALCHEMYAPI}`,
    accounts: [`0x${process.env.ACCOUNT_PRIVATE_KEY}`],
  },
  rinkeby: {
    url: `${process.env.RINKEBY_ALCHEMYAPI}`,
    accounts: [`0x${process.env.ACCOUNT_PRIVATE_KEY}`],
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
