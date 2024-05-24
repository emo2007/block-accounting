//SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@chainlink/contracts/src/v0.8/ChainlinkClient.sol";
import "@chainlink/contracts/src/v0.8/shared/access/ConfirmedOwner.sol";

contract Agreement is ChainlinkClient, ConfirmedOwner {
    using Chainlink for Chainlink.Request;

    address private oracleAddress;
    bytes32 private jobId;
    uint256 private fee;
    address public multisigWallet;

    constructor(
        address _chainLinkToken,
        address _oracleAddress,
        string memory _jobId,
        uint _fee,
        address _multiSigAddress
    ) ConfirmedOwner(_multiSigAddress) {

        _setChainlinkToken(_chainLinkToken);

        setOracleAddress(_oracleAddress);

        setJobId(_jobId);

        setFeeInHundredthsOfLink(_fee);

        multisigWallet = _multiSigAddress;
    }


    // Send a request to the Chainlink oracle
    function request(string memory url) public {

        Chainlink.Request memory req = _buildOperatorRequest(jobId, this.fulfill.selector);

        req._add('method', 'GET');
        req._add('url', url);
        req._add('headers', '["content-type", "application/json"]');
        req._add('body', '');
        req._add('contact', '');
        req._add('path', '');
        _sendOperatorRequest(req, fee);
    }

    bool public response;

    // Receive the result from the Chainlink oracle
    event RequestFulfilled(bytes32 indexed requestId);
    function fulfill(bytes32 requestId, bool data) public recordChainlinkFulfillment(requestId) {
        emit RequestFulfilled(requestId);
        response = data;
    }

    // Update oracle address
    function setOracleAddress(address _oracleAddress) public onlyOwner {
        oracleAddress = _oracleAddress;
        _setChainlinkOracle(_oracleAddress);
    }
    function getOracleAddress() public view onlyOwner returns (address) {
        return oracleAddress;
    }

    // Update jobId
    function setJobId(string memory _jobId) public onlyOwner {
        jobId = bytes32(bytes(_jobId));
    }
    function getJobId() public view onlyOwner returns (string memory) {
        return string(abi.encodePacked(jobId));
    }

    // Update fees
    function setFeeInJuels(uint256 _feeInJuels) public onlyOwner {
        fee = _feeInJuels;
    }
    function setFeeInHundredthsOfLink(uint256 _feeInHundredthsOfLink) public onlyOwner {
        setFeeInJuels((_feeInHundredthsOfLink * LINK_DIVISIBILITY) / 100);
    }
    function getFeeInHundredthsOfLink() public view onlyOwner returns (uint256) {
        return (fee * 100) / LINK_DIVISIBILITY;
    }

    function withdrawLink() public onlyOwner {
        LinkTokenInterface link = LinkTokenInterface(_chainlinkTokenAddress());
        require(
            link.transfer(msg.sender, link.balanceOf(address(this))),
            "Unable to transfer"
        );
    }
}