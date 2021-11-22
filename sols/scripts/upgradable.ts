import { ethers, upgrades } from "hardhat";

const code = "TheQTokV2";

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log(`Deploy contracts with the account: ${deployer.address}`);
  const balance = await deployer.getBalance();

  console.log(`Account balance: ${balance.toString()}`);

  const Token = await ethers.getContractFactory(code);
  const token = await Token.deploy();

  await token.deployed();

  console.log(`${code} deployed to:${token.address}`);
}

async function update() {
  const addr = "0xE4f53b165cCbC3F81B3813a4aD4d505ae8898ADA";
  console.log(`Upgrade for ${addr}`);
  const TokenUpdated = await ethers.getContractFactory(code);
  const upgraded = await upgrades.upgradeProxy(addr, TokenUpdated);
  await upgraded.deployed();
  console.log(`The ${code} is upgrade at ${upgraded.address}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
