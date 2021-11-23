//SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "hardhat/console.sol";

contract Greeter {
    string private value;
    Person[] private people;

    mapping(uint256 => Person) public peopleMap;
    uint256 private peopleCount;

    struct Person {
        string _firstName;
        string _lastName;
    }

    address owner;

    modifier onlyMySelf() {
        require(msg.sender == owner);
        _;
    }

    uint256 openTime = 1637633801 + 5 * 60;

    modifier onlyWhileOpen() {
        require(block.timestamp >= openTime);
        _;
    }

    constructor() {
        value = "hello world";
        // set owner to the one who deploys this contract
        owner = msg.sender;
        console.log("Initialize the system");
    }

    function get() public view returns (string memory) {
        return value;
    }

    function set(string memory _value) public {
        value = _value;
    }

    function addPerson(string memory _first, string memory _last)
        public
        onlyMySelf
    {
        peopleCount += 1;
        Person memory p = Person(_first, _last);
        people.push(p);
        peopleMap[peopleCount] = p;
    }

    function peoples() public view returns (Person[] memory) {
        return people;
    }

    function totalPeople() public view returns (uint256) {
        return people.length;
    }
}
