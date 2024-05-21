//SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@chainlink/contracts/src/v0.8/ChainlinkClient.sol";
import "@chainlink/contracts/src/v0.8/shared/access/ConfirmedOwner.sol";

/**
 * Request testnet LINK and ETH here: https://faucets.chain.link/
 * Find information on LINK Token Contracts and get the latest ETH and LINK faucets here: https://docs.chain.link/docs/link-token-contracts/
 */

/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 */

contract LinkWellBoolConsumerContractExample is ChainlinkClient, ConfirmedOwner {
    using Chainlink for Chainlink.Request;

    address private oracleAddress;
    bytes32 private jobId;
    uint256 private fee;

    constructor(
        address _chainLinkToken,
        address _oracleAddress,
        string memory _jobId,
        uint _fee,
        address _multiSigAddress,
        address[] memory _owners,
        uint[] memory _shares
    ) ConfirmedOwner(_multiSigAddress) {

        _setChainlinkToken(_chainLinkToken);

        setOracleAddress(_oracleAddress);

        setJobId(_jobId);

        setFeeInHundredthsOfLink(_fee);

        multisigWallet = _multiSigAddress;
    }

    constructor() ConfirmedOwner(msg.sender) {
        _setChainlinkToken(0x0Fd9e8d3aF1aaee056EB9e802c3A762a667b1904);
        setOracleAddress(0xd36c6B1777c7f3Db1B3201bDD87081A9045B7b46);
        setJobId("43309009a154495cb2ed794233e6ff56");
        setFeeInHundredthsOfLink(0);     // 0 LINK
    }

    // Send a request to the Chainlink oracle
    function request() public {

        Chainlink.Request memory req = _buildOperatorRequest(jobId, this.fulfill.selector);

        // DEFINE THE REQUEST PARAMETERS (example)
        req._add('method', 'POST');
        req._add('url', 'https://httpbin.org/post');
        req._add('headers', '["accept", "application/json", "set-cookie", "sid=14A52"]');
        req._add('body', '{"data":[{"coin":"BTC","isActive":false},{"coin":"ETH","isActive":false},{"coin":"LINK","isActive":true}]}');
        req._add('contact', '');     // PLEASE ENTER YOUR CONTACT INFO. this allows us to notify you in the event of any emergencies related to your request (ie, bugs, downtime, etc.). example values: 'derek_linkwellnodes.io' (Discord handle) OR 'derek@linkwellnodes.io' OR '+1-617-545-4721'

        // The following curl command simulates the above request parameters:
        // curl 'https://httpbin.org/post' --request 'POST' --header 'content-type: application/json' --header 'set-cookie: sid=14A52' --data '{"data":[{"coin":"BTC","isActive":false},{"coin":"ETH","isActive":false},{"coin":"LINK","isActive":true}]}'

        // PROCESS THE RESULT (example)
        req._add('path', 'json,data,2,isActive');

        // Send the request to the Chainlink oracle
        _sendOperatorRequest(req, fee);
    }

    bool public response;

    // Receive the result from the Chainlink oracle
    event RequestFulfilled(bytes32 indexed requestId);
    function fulfill(bytes32 requestId, bool data) public recordChainlinkFulfillment(requestId) {
        // Process the oracle response
        // emit RequestFulfilled(requestId);    // (optional) emits this event in the on-chain transaction logs, allowing Web3 applications to listen for this transaction
        response = data;     // example value: true
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