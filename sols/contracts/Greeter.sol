//SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "hardhat/console.sol";
import "hardhat/console.sol";

contract Greeter {
    string private value;

    Person[] private people;

    struct Person {
        string _firstName;
        string _lastName;
    }

    constructor() {
        value = "hello world";
    }

    function get() public view returns (string memory) {
        return value;
    }

    function set(string memory _value) public {
        value = _value;
    }

    function addPerson(string memory _first, string memory _last) public {
        people.push(Person(_first, _last));
    }

    function peoples() public view returns (Person[] memory) {
        return people;
    }

    function totalPeople() public view returns (uint256) {
        return people.length;
    }
}
