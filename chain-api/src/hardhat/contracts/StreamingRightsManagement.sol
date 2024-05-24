//SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@chainlink/contracts/src/v0.8/ChainlinkClient.sol";
import "@chainlink/contracts/src/v0.8/shared/access/ConfirmedOwner.sol";
import "./Payroll.sol";


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
        uint[] memory _shares
    ) ConfirmedOwner(_multiSigAddress) {

        _setChainlinkToken(_chainLinkToken);

        setOracleAddress(_oracleAddress);

        setJobId(_jobId);

        setFeeInHundredthsOfLink(_fee);

        multisigWallet = _multiSigAddress;


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
    modifier hasValidPayoutContract() {
        require(address(payoutContract) != address(0), "payoutContract not initialized");
        _;
    }

    function getShare(address owner) public view returns(uint){
        return ownerShare[owner];
    }

    function setPayoutContract(address payable _payoutAddress) public onlyOwner {
        require(_payoutAddress != address(0), "Invalid address: zero address not allowed");
        payoutContract = Payroll(_payoutAddress);
    }

    // Send a request to the Chainlink oracle
    function request(string memory url) external onlyOwner{

        Chainlink.Request memory req = _buildOperatorRequest(jobId, this.fulfill.selector);

        req._add('method', 'GET');
        req._add('url', url);

        //if returns just int - then empty path
        req._add('path', '');
        req._addInt('multiplier', 10 ** 18);
        req._add('headers', '["content-type", "application/json"]');
        req._add('body', '');
        req._add('contact', '');
        _sendOperatorRequest(req, fee);
    }

    uint256 public totalPayoutInUSD;

    // Receive the result from the Chainlink oracle
    event RequestFulfilled(bytes32 indexed requestId);

    function fulfill(bytes32 requestId, uint256 data) public recordChainlinkFulfillment(requestId) {
        totalPayoutInUSD = data / 1e18;
    }

    function payout() external onlyOwner hasValidPayoutContract{

        // using arrays to reduce gas
        uint[] memory shares = new uint[](owners.length);

    for(uint i=0; i< owners.length; i++){
          shares[i] = ownerShare[owners[i]] * totalPayoutInUSD / 100;
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