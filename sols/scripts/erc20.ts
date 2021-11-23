import { deploy } from '@openzeppelin/hardhat-upgrades/dist/utils';
import { ethers } from 'hardhat';

async function main() {
    const [deployer, addr1] = await ethers.getSigners();
    const addr = deployer.address;
    console.log(`Deploy using accounts ${addr} with balances ${await deployer.getBalance()}`);

    const Tok = await ethers.getContractFactory("MyQToken");
    const tok = await Tok.deploy();
    console.log("Wait token for deploying");
    await tok.deployed();
    console.log(`Got token address ${tok.address}, deployed with wallet ${addr1.address}`);
}
// 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0


main();