import { ethers } from "hardhat";

async function main() {
	const Contract = await ethers.getContractFactory("Greeter");
	const c = await Contract.deploy();

	await c.deployed();
	console.log(`contract deployed at ${c.address}`)


}

async function interact() {
	const addr = '0x5FbDB2315678afecb367f032d93F642f64180aa3';
	const Contract = await ethers.getContractFactory("Greeter");
	// let interacts with tokens
	const attach = await Contract.attach(addr);

	const before = await attach.get();
	console.log("Before settings", before);
	// let's set some value

	await attach.set(before + " my name is q 2");
	console.log(`after setup: ${await attach.get()}`);
}

interact();