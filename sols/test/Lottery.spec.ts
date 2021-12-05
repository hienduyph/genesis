import { ethers, network } from "hardhat";
import net from "../lib/network";

describe("Lottery Testings", () => {
  it("Deploy full steps", async () => {
    const [owner] = await ethers.getSigners();
    const conf = net[network.name];

    console.log(`Deploy using ${owner.address} with balance ${await owner.getBalance()}`)

    const Token = await ethers.getContractFactory("Lottery");
    const tok = await Token.deploy(
      conf.feedAddr,
      conf._vrfCoordinator,
      conf._link,
      conf._fee,
      conf._keyHash,
    );
    await tok.deployed();

    console.log(`Token deploy at ${tok.address}`);
  });
});
