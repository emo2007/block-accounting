//SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@chainlink/contracts/src/v0.8/ChainlinkClient.sol";
import "@chainlink/contracts/src/v0.8/shared/access/ConfirmedOwner.sol";
import "./Payroll.sol";
/**
 * Request testnet LINK and ETH here: https://faucets.chain.link/
 * Find information on LINK Token Contracts and get the latest ETH and LINK faucets here: https://docs.chain.link/docs/link-token-contracts/
 */

/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 */

contract StreamingRightsManagement is ChainlinkClient, ConfirmedOwner {
    using Chainlink for Chainlink.Request;

    address public oracleAddress;
    bytes32 private jobId;
    uint256 private fee;
    address public multisigWallet;

    mapping(address => uint) public ownerShare;
    address[] public owners;

    Payroll public payoutContract;

    constructor(
        address _chainLinkToken,
        address _oracleAddress,
        string memory _jobId,
        uint _fee,
        address _multiSigAddress,
        address[] memory _owners,
        uint[] memory _shares,
        address payable _payoutAddress
    ) ConfirmedOwner(_multiSigAddress) {

        _setChainlinkToken(_chainLinkToken);

        setOracleAddress(_oracleAddress);

        setJobId(_jobId);

        setFeeInHundredthsOfLink(_fee);

        multisigWallet = _multiSigAddress;

        payoutContract = Payroll(_payoutAddress);

        require(_owners.length == _shares.length, "Owners and shares length mismatch");

        uint sumShare = 0;

        for(uint i=0; i<_shares.length;i++){
            sumShare += _shares[i];
        }

        require(sumShare ==100, 'Invalid share percentage');
        for (uint i = 0; i < _owners.length; i++) {
            require(_shares[i] > 0, 'Share cannot be less than 0');
            ownerShare[_owners[i]] = _shares[i];
            owners.push(_owners[i]);
        }
    }
    //get share
    //update share
    //change payout address
    //
    function getShare(address owner) public returns(uint){
        return ownerShare[owner];
    }


    // Send a request to the Chainlink oracle
    function request() public {

        Chainlink.Request memory req = _buildOperatorRequest(jobId, this.fulfill.selector);

        // DEFINE THE REQUEST PARAMETERS (example)
        req._add('method', 'GET');
        req._add('url', 'https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC,ETH&tsyms=USD,EUR');
        req._add('headers', '["content-type", "application/json", "set-cookie", "sid=14A52"]');
        req._add('body', '');
        req._add('contact', '');     // PLEASE ENTER YOUR CONTACT INFO. this allows us to notify you in the event of any emergencies related to your request (ie, bugs, downtime, etc.). example values: 'derek_linkwellnodes.io' (Discord handle) OR 'derek@linkwellnodes.io' OR '+1-617-545-4721'

        // The following curl command simulates the above request parameters:
        // curl 'https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC,ETH&tsyms=USD,EUR' --request 'GET' --header 'content-type: application/json' --header 'set-cookie: sid=14A52'

        // PROCESS THE RESULT (example)
        req._add('path', 'ETH,USD');
        // Send the request to the Chainlink oracle
        _sendOperatorRequest(req, fee);
    }

    uint256 public totalPayoutInUSD;

    // Receive the result from the Chainlink oracle
    event RequestFulfilled(bytes32 indexed requestId);

    function fulfill(bytes32 requestId, uint256 data) public recordChainlinkFulfillment(requestId) {
        // Process the oracle response
        // emit RequestFulfilled(requestId);    // (optional) emits this event in the on-chain transaction logs, allowing Web3 applications to listen for this transaction
        totalPayoutInUSD = data / 100;     // example value: 1875870000000000000000 (1875.87 before "multiplier" is applied)
    }

    function payout() external onlyOwner {
        // using arrays to reduce gas
        uint[] memory shares;

        for(uint i=0; i< owners.length; i++){
          shares[i] = ownerShare[owners[i]];
        }
        payoutContract.oneTimePayout(owners, shares);
    }

    // Update oracle address
    function setOracleAddress(address _oracleAddress) public onlyOwner {
        oracleAddress = _oracleAddress;
        _setChainlinkOracle(_oracleAddress);
    }

    function getOracleAddress() public view returns (address) {
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