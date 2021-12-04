// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
import "@chainlink/contracts/src/v0.8/VRFConsumerBase.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "hardhat/console.sol";

contract Lottery is Ownable, VRFConsumerBase {
    address[] public players;
    address public winner;
    uint256 public randomness;

    uint256 public usdEntryFee;
    AggregatorV3Interface public usdPriceFeed;
    enum LOTTERY_STATE {
        OPEN,
        CLOSED,
        CALCULATING_WINNER
    }
    LOTTERY_STATE public state;
    uint256 public fee;
    bytes32 public keyHash;

    constructor(
        address feedAddr,
        address _vrfCoordinator,
        address _link,
        uint256 _fee,
        bytes32 _keyHash
    ) VRFConsumerBase(_vrfCoordinator, _link) {
        usdEntryFee = 50 * 10**18;
        usdPriceFeed = AggregatorV3Interface(feedAddr);
        state = LOTTERY_STATE.CLOSED;
        keyHash = _keyHash;
        fee = _fee;
    }

    function enter() public payable {
        require(state == LOTTERY_STATE.OPEN, "Lottery is not open!");
        require(
            getConversionRate(msg.value) >= usdEntryFee,
            "You must send more ETH!"
        );
        players.push(msg.sender);
    }

    function getConversionRate(uint256 ethaMount)
        public
        view
        returns (uint256)
    {
        (, int256 price, , , ) = usdPriceFeed.latestRoundData();
        uint256 factor = 18 - usdPriceFeed.decimals();
        uint256 rate = uint256(price) * 10**factor;
        return (rate * ethaMount) / 10**18;
    }

    function start() public onlyOwner {
        require(state == LOTTERY_STATE.CLOSED, "Can not start a new lotter");
        state = LOTTERY_STATE.OPEN;
    }

    function end() public onlyOwner {
        require(state == LOTTERY_STATE.OPEN, "Lottery is not open yet!");

        state = LOTTERY_STATE.CALCULATING_WINNER;
        requestRandomness(keyHash, fee);
    }

    /**
     * Callback function used by VRF Coordinator
     */
    function fulfillRandomness(bytes32 _v, uint256 _randomness)
        internal
        override
    {
        require(
            state == LOTTERY_STATE.CALCULATING_WINNER,
            "You aren't there yet!"
        );
        require(_randomness > 0, "random must > 0");

        randomness = _randomness;

        winner = players[randomness % players.length];

        uint256 total = address(this).balance;
        console.log("Winner Is %s, total %s", winner, total);
        payable(winner).transfer(total);

        // reset data
        players = new address[](0);
    }
}
