import { ethers, network } from "hardhat";

import { sleep } from "../lib/core";
import { transfer } from "../lib/link";
import net from "../lib/network";
import { ask } from "../lib/prompt";

const conf = net[network.name];

async function main() {
  const [owner] = await ethers.getSigners();
  console.log(`Deploy using ${owner.address} with balance ${await owner.getBalance()}`)


  const Token = await ethers.getContractFactory("Lottery");

  if ("y" !== await ask("Are you ready (y/n)")) {
    console.log("Good bye!");
    return
  }
  console.log("Deploying the Lotter");
  const tok = await Token.deploy(
    conf.feedAddr,
    conf._vrfCoordinator,
    conf._link,
    ethers.utils.parseEther(conf._fee),
    conf._keyHash,
  );
  await tok.deployed();

  console.log(`Token deploy at ${tok.address}`);
  // let's transfer some LINK for the contract addr
  await transfer(conf._link, tok.address, conf._fee);

  play(tok.address);
}

async function play(addr: string) {
  const accs = await ethers.getSigners();
  const Contract = await ethers.getContractFactory("Lottery");
  const lottery = Contract.attach(addr);

  console.log("Start the lottery");
  await lottery.start();

  while (await lottery.state() !== 0) {
    console.log("Wait for lottery active...");
    await sleep(2);
  }

  for (const acc of accs) {
    console.log(`Enters ${acc.address}`)
    await lottery.connect(acc).enter({ value: ethers.utils.parseEther("0.014") });
  }

  console.log("End the lottery");
  await lottery.end();

  console.log("Waiting the lottery to calc winner");

  while (true) {
    await sleep(10);

    console.log("Checking calc done!");
    const w = await lottery.winner();
    if (w.toString() === "0x0000000000000000000000000000000000000000") {
      continue
    }
    console.log(`Winner is ${w}`);
    return;
  }
}

main();