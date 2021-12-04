// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./SimpleStorage.sol";

contract StorageFactory is SimpleStorage {
    SimpleStorage[] public storages;

    function createSimpleStorageContract() public {
        SimpleStorage ss = new SimpleStorage();
        storages.push(ss);
    }

    function fStore(uint256 idx, uint256 value) public {
        // get the ABI for adress at index `idx`
        SimpleStorage ss = SimpleStorage(address(storages[idx]));
        ss.store(value);
    }

    function fRetrive(uint256 idx) public view returns (uint256) {
        return SimpleStorage(address(storages[idx])).retrive();
    }
}
