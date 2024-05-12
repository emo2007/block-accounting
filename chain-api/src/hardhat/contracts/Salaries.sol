// SPDX-License-Identifier: MIT

pragma solidity ^0.8.7;

import {AggregatorV3Interface} from '@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol';

contract Salaries {
    AggregatorV3Interface internal dataFeed;
    address public multisigWallet;
    mapping(address => uint) public salaries;

    //0xF0d50568e3A7e8259E16663972b11910F89BD8e7
    constructor(address _multisigWallet, address _priceFeedAddress) {
        multisigWallet = _multisigWallet;
        dataFeed = AggregatorV3Interface(_priceFeedAddress);
    }

    modifier onlyMultisig() {
        require(msg.sender == multisigWallet, 'Unauthorized');
        _;
    }

    function getSalary(address employee) public view returns(uint) {
        return salaries[employee];
    }

    function getLatestUSDTPriceInETH() public view returns (int) {
        (
            ,
            /* uint80 roundID */ int answer /* uint startedAt */ /* uint timeStamp */ /* uint80 answeredInRound */,
            ,
            ,

        ) = dataFeed.latestRoundData();
        return answer;
    }

    function setSalary(
        address employee,
        uint salaryInUSDT
    ) external onlyMultisig {
        salaries[employee] = salaryInUSDT;
    }

    function payoutInETH(address employee) external onlyMultisig {
        uint salaryInUSDT = salaries[employee];
        require(salaryInUSDT > 0, 'No salary set');

        int ethToUSDT = getLatestUSDTPriceInETH();
        require(ethToUSDT > 0, 'Invalid price data');

        // Convert salary from USDT to ETH based on the latest price
        uint salaryInETH = uint(salaryInUSDT * 1e18) / uint(ethToUSDT);

        // Check sufficient balance
        require(
            address(this).balance >= salaryInETH,
            'Insufficient contract balance'
        );

        salaries[employee] = 0; // Reset salary after payment
        payable(employee).transfer(salaryInETH);
    }

    // Fallback to receive ETH
    receive() external payable {}
}
