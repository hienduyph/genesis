// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Math.sol";

contract ERC20Token {
    string public name;
    mapping(address => uint256) public balances;

    constructor(string memory _name) {
        name = _name;
    }

    function mint() public virtual {
        balances[msg.sender] += 1;
    }
}

contract MyQToken is ERC20Token {
    using Math for uint256;

    string public symbol = "MyQToken";
    address[] public owners;
    uint256 ownerCount;

    uint256 public value;

    constructor() ERC20Token("MyQToken") {}

    function mint() public override {
        super.mint();
        ownerCount++;
        owners.push(msg.sender);
    }

    function calc(uint256 a, uint256 b) public {
        value = a.divide(b);
    }
}
