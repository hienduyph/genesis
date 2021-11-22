import { ethers } from "hardhat";

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log(`Deploy contracts with the account: ${deployer.address}`);
  const balance = await deployer.getBalance();

  console.log(`Account balance: ${balance.toString()}`);

  const Token = await ethers.getContractFactory("TheQTokV2");
  const token = await Token.deploy();

  await token.deployed();

  console.log("Greeter deployed to:", token.address);
}

main()
  .then(() => process.exit(0))
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
