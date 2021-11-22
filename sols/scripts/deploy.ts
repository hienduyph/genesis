import { ethers } from "hardhat";

async function main() {
  const code = "TheQTokV2";
  const [deployer] = await ethers.getSigners();
  console.log(`Deploy contracts with the account: ${deployer.address}`);
  const balance = await deployer.getBalance();

  console.log(`Account balance: ${balance.toString()}`);

  const Token = await ethers.getContractFactory(code);
  const token = await Token.deploy();

  await token.deployed();

  console.log(`${code} deployed to:${token.address}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
