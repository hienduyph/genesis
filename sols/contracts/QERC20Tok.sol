// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

contract QERC20Tok {
    string public name = "Jet Another Tok";
    string public symbol = "QERC20Tok";

    mapping(address => uint256) balances;

    address payable wallet;

    // add indexed coult filter events by buyter
    // we could listen event and do something more than things
    event Purchase(address indexed _buyer, uint256 _amount);

    constructor(address payable _wallet) {
        wallet = _wallet;
    }

    fallback() external payable {
        buyToken();
    }

    function buyToken() public payable {
        balances[msg.sender] += 1;
        wallet.transfer(msg.value);
        emit Purchase(msg.sender, 1);
    }

    function balanceOf(address _acc) external view returns (uint256) {
        return balances[_acc];
    }
}
