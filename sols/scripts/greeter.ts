import { ethers } from "hardhat";

async function deploy() {
	const [deployer] = await ethers.getSigners();
	const balance = await deployer.getBalance();
	const account = await deployer.getAddress();
	console.log(`Deploy with account ${account} and balance ${balance}`);

	const Contract = await ethers.getContractFactory("Greeter");
	const c = await Contract.deploy();

	await c.deployed();
	console.log(`contract deployed at ${c.address}`)
}

async function interact() {
	const [deployer, otherAddr] = await ethers.getSigners();

	const addr = process.env.CONTRACT_ADDR || '';
	const Contract = await ethers.getContractFactory("Greeter");
	// let interacts with tokens
	const cont = await Contract.attach(addr);

	const before = await cont.get();
	console.log("Before settings", before);
	// let's set some value

	await cont.set(before + " my name is q 2");
	console.log(`after setup: ${await cont.get()}`);

	// test the modifier
	await cont.connect(otherAddr.address).addPerson("mxyh", "2");
}

deploy();