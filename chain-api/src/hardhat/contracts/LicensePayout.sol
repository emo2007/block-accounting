// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import {AggregatorV3Interface} from '@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol';

contract LicensePayout {
    AggregatorV3Interface internal dataFeed;
    address public multisigWallet;
    mapping(address => uint) public salaries;
    event Payout(address indexed employee, uint salaryInETH);
    event PayoutFailed(address indexed employee, uint salaryInETH, string reason);
    constructor(address _multisigWallet, address _priceFeedAddress) {
        multisigWallet = _multisigWallet;
        dataFeed = AggregatorV3Interface(_priceFeedAddress);
    }

    modifier onlyMultisig() {
        require(msg.sender == multisigWallet, 'Unauthorized');
        _;
    }

    function getUsdtSalary(address employee) public view returns(uint) {
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

    function oneTimePayout(address payable employee) external onlyMultisig {

    }

    function setSalary(
        address employee,
        uint salaryInUSDT
    ) external onlyMultisig {
        salaries[employee] = salaryInUSDT;
    }

    function getEmployeeSalaryInEth(address employee) public view returns(uint){
        uint salaryInUSDT = salaries[employee];
        require(salaryInUSDT > 0, 'No salary set');

        int ethToUSDT = getLatestUSDTPriceInETH();
        require(ethToUSDT > 0, 'Invalid price data');
        uint salaryInETH = uint(salaryInUSDT * 1e18) / uint(ethToUSDT);
        return salaryInETH * 1e8;
    }

    function payoutInETH(address payable employee) external onlyMultisig {
        uint salaryInETH = getEmployeeSalaryInEth(employee);
        // Check sufficient balance
        require(
            address(this).balance >= salaryInETH,
            'Insufficient contract balance'
        );

        (bool success, ) = employee.call{value: salaryInETH}("");
        if (success) {
            emit Payout(employee, salaryInETH);
        } else {
            emit PayoutFailed(employee, salaryInETH, "Transfer failed");
        }
    }

    // Fallback to receive ETH
    receive() external payable {}
}
