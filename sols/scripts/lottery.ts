import { ethers } from "hardhat";
import conf from "../lib/kovan";
import { ask } from "../lib/prompt";

async function main() {
  const [owner] = await ethers.getSigners();
  console.log(`Deploy using ${owner.address} with balance ${await owner.getBalance()}`)

  const Token = await ethers.getContractFactory("Lottery");

  if ("y" !== await ask("Are you ready")) {
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
  console.log(`Please send mimimun ${conf._fee} LINK to ${tok.address}`);
  const v = await ask("Enter OK after done")
  if (v !== "OK") {
    console.log("Skip playing with the contract");
    return
  }
  play(tok.address);
}

async function play(addr: string) {
  const [owner, acc1, acc2] = await ethers.getSigners();
  const Contract = await ethers.getContractFactory("Lottery");
  const lottery = Contract.attach(addr);

  console.log("Start the lottery");
  await lottery.start();

  // now let's enter
  console.log("Owner enters")
  await lottery.enter({ value: ethers.utils.parseEther("0.014") });

  console.log(`Enters ${acc1.address}`);
  await lottery.connect(acc1).enter({ value: ethers.utils.parseEther("0.013") });

  console.log(`Enters ${acc2.address}}`);
  await lottery.connect(acc2).enter({ value: ethers.utils.parseEther("0.013") });

  console.log("End the lottery");
  await lottery.end();

  console.log("Waiting the lottery to calc winner");

  while (true) {
    await new Promise(r => setTimeout(r, 10000));

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