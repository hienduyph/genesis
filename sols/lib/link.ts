import abi from "@chainlink/contracts/abi/v0.4/LinkToken.json";
import { ethers } from "hardhat";

export const transfer = async (linkAddr: string, contractAddr: string, amount: string) => {
  console.log(`Transfer ${amount} LINK at ${linkAddr} for ${contractAddr}`);
  const [owner] = await ethers.getSigners();
  const linkToken = new ethers.Contract(linkAddr, abi, ethers.provider);
  await linkToken.connect(owner).transfer(contractAddr, ethers.utils.parseEther(amount))
};