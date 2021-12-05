const rinkeby = {
  feedAddr: "0x8A753747A1Fa494EC906cE90E9f37563A8AF630e",
  _vrfCoordinator: "0xb3dCcb4Cf7a26f6cf6B120Cf5A73875B7BBc655B",
  _link: "0x01BE23585060835E02B77ef475b0Cc51aA1e0709",
  _fee: "0.1",
  _keyHash: "0x2ed0feb3e7fd2022120aa84fab1945545a9f2ffc9076fd6156fa96eaff4c1311",
}

const kovan = {
  feedAddr: "0x9326BFA02ADD2366b30bacB125260Af641031331",
  _vrfCoordinator: "0xdD3782915140c8f3b190B5D67eAc6dc5760C46E9",
  _link: "0xa36085F69e2889c224210F603D836748e7dC0088",
  _fee: "0.1",
  _keyHash: "0x6c3699283bda56ad74f6b855546325b68d482e983852a7a82979cc4807b641f4",
}

interface NetConf {
  feedAddr: string;
  _vrfCoordinator: string;
  _link: string;
  _fee: string;
  _keyHash: string;
}

const c: { [key: string]: NetConf } = { kovan, rinkeby };

export default c;