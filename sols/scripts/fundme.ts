import { ethers } from "hardhat";

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log(`Deploy using addr ${deployer.address} and balance ${await deployer.getBalance()}`)

    const Tok = await ethers.getContractFactory("FundMe");
    const tok = await Tok.deploy();

    await tok.deployed();
    console.log("Tok deploy at", tok.address);
}

main();