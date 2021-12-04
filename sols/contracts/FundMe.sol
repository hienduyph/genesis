// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";

contract FundMe {
    mapping(address => uint256) public addressToAmountFunded;
    address[] public funders;
    address public owner;

    // Rinkeby test net
    address private feedAddr = 0x8A753747A1Fa494EC906cE90E9f37563A8AF630e;

    constructor() {
        owner = msg.sender;
    }

    function fund() public payable {
        uint256 minUSD = 50 * 10**18; // 50USD
        uint256 fundedUSD = getConversionRate(msg.value);
        require(fundedUSD >= minUSD, "You need to spend more ETH!");

        addressToAmountFunded[msg.sender] += msg.value;
        funders.push(msg.sender);
    }

    function withdraw() public payable onlyOwner {
        payable(msg.sender).transfer(address(this).balance);
        // reset data
        for (uint256 i = 0; i < funders.length; i++) {
            addressToAmountFunded[funders[i]] = 0;
        }
        funders = new address[](0);
    }

    function getConversionRate(uint256 ethAmount)
        public
        view
        returns (uint256)
    {
        uint256 price = getPrice();
        uint256 ethPriceInUSD = price * ethAmount;
        return ethPriceInUSD / 10**18;
    }

    /**
     *  Get real price in usd (in wei)
     */
    function getPrice() public view returns (uint256) {
        AggregatorV3Interface feed = AggregatorV3Interface(feedAddr);
        (, int256 price, , , ) = feed.latestRoundData();
        uint256 factor = 18 - feed.decimals();
        return uint256(price) * 10**factor;
    }

    function getVersion() public view returns (uint256) {
        AggregatorV3Interface feed = AggregatorV3Interface(feedAddr);
        return feed.version();
    }

    modifier onlyOwner() {
        require(owner == msg.sender, "Only owner can withdraw");
        _;
    }
}
