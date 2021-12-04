// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

contract SimpleStorage {
    uint256 private favoriteNumber;

    function store(uint256 _num) public {
        favoriteNumber = _num;
    }

    function retrive() public view returns (uint256) {
        return favoriteNumber;
    }
}
