import { ethers } from "hardhat";
import ProxyAdminABI from "@openzeppelin/contracts/build/contracts/ProxyAdmin.json";
import TransparentUpgradeableProxyABi from "@openzeppelin/contracts/build/contracts/TransparentUpgradeableProxy.json";
import { sleep } from "../lib/core";

/**
 * Rules:
 * - has initializer:
 *  + intializer.encode_input()
 *  + proxyAdmin.upgradeAndCall()
 *  + proxy.upgradeToAndCall()
 * - else:
 *  + fixed; `0x`
 *  + proxyAdmin.upgrade()
 *  + proxy.upgradeTo()
 */
async function main() {
  const [deployer] = await ethers.getSigners();
  const ProxyAdmin = new ethers.ContractFactory(ProxyAdminABI.abi, ProxyAdminABI.bytecode, deployer);
  const UpgradableProxy = new ethers.ContractFactory(TransparentUpgradeableProxyABi.abi, TransparentUpgradeableProxyABi.bytecode, deployer);

  const BoxV1 = await ethers.getContractFactory("Box");
  const BoxV2 = await ethers.getContractFactory("BoxV2");

  // let's deploy some things
  const proxyAdmin = await ProxyAdmin.deploy();
  const box = await BoxV1.deploy();

  console.log(`Deployed  proxyAdmin ${proxyAdmin.address}; box ${box.address}`);

  await sleep(10);

  const proxy = await UpgradableProxy.deploy(box.address, proxyAdmin.address, "0x")
  console.log(`proxy Deployed ${proxy.address}`);


  // from now we can using proxy address to call the function
  let proxyBox = BoxV1.attach(proxy.address);
  await proxyBox.setValue(1111);

  await sleep(20);

  console.log("Value now is ", (await proxyBox.retrieve()).toString());

  // now let upgrade to v2
  const boxV2 = await BoxV2.deploy();
  await sleep(10);
  console.log(`BoxV2 deployed at ${boxV2.address}`)
  await proxyAdmin.upgrade(proxy.address, boxV2.address);


  await sleep(10);
  // let usign new v2 impl
  proxyBox = BoxV2.attach(proxy.address);

  await proxyBox.increment();
  await sleep(10);
  console.log("Value after increment is ", (await proxyBox.retrieve()).toString());

}

main();
