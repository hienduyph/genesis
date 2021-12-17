// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";

contract NFTSimple is ERC721URIStorage {
    uint256 internal tokenCounter;

    constructor() ERC721("NFTSimple", "QNFTS") {}

    function createCollective(string memory tokenURI) public returns (uint256) {
        uint256 newTokenID = tokenCounter;
        _safeMint(msg.sender, newTokenID);
        _setTokenURI(newTokenID, tokenURI);
        tokenCounter += 1;
        return newTokenID;
    }
}
