import { ethers } from "hardhat";
import { OPENSEA_URI } from "../lib/constants";

const sample_token_uri = "https://ipfs.io/ipfs/Qmd9MCGtdVz2miNumBHDbvj8bigSgTwnr4SbyH6DNnpWdt?filename=0-PUG.json"


async function main() {
  const Contract = await ethers.getContractFactory("NFTSimple");
  const token = await Contract.deploy();
  console.log(`Contract deployed at ${token.address}`);

  const tokenID = await token.createCollective(sample_token_uri);
  console.log(
    `Awesome, you can view your NFT at ${OPENSEA_URI(token.address, tokenID)}`
  )
  console.log("Please wait up to 20 minutes, and hit the refresh metadata button");
}

main();