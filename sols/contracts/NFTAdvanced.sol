// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@chainlink/contracts/src/v0.8/VRFConsumerBase.sol";

contract NFTAdvanced is ERC721URIStorage, VRFConsumerBase {
    uint256 public tokenCounter;
    bytes32 public keyHash;
    uint256 public fee;

    enum Breed {
        PUG,
        SHIBA_INU,
        ST_BERNARD
    }
    mapping(uint256 => Breed) tokenIDtoBreed;
    mapping(bytes32 => address) public requestIDToSender;

    event requestedCollectible(bytes32 requestid, address req);
    event breedAssigned(uint256 indexed tokenId, Breed breed);

    constructor(
        address _vrfCoodrinator,
        address _link,
        bytes32 _keyHash,
        uint256 _fee
    )
        VRFConsumerBase(_vrfCoodrinator, _link)
        ERC721("NFTAdvanced", "QNFTAdvanced")
    {
        keyHash = _keyHash;
        fee = _fee;
    }

    function createCollectible() public returns (bytes32) {
        bytes32 requestId = requestRandomness(keyHash, fee);
        requestIDToSender[requestId] = msg.sender;
        emit requestedCollectible(requestId, msg.sender);
        return requestId;
    }

    function fulfillRandomness(bytes32 requestID, uint256 randomNumber)
        internal
        override
    {
        Breed bree = Breed(randomNumber % 3);
        uint256 newTokenId = tokenCounter;
        tokenIDtoBreed[newTokenId] = bree;
        _safeMint(requestIDToSender[requestID], newTokenId);
        emit breedAssigned(newTokenId, bree);

        tokenCounter += 1;
    }

    function setTokenURI(uint256 tokenId, string memory tokenURI) public {
        require(
            _isApprovedOrOwner(_msgSender(), tokenId),
            "caller is not owner no approved"
        );
        _setTokenURI(tokenId, tokenURI);
    }
}
