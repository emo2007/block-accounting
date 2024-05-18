// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import {AggregatorV3Interface} from '@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol';

contract Payroll {
    AggregatorV3Interface internal dataFeed;
    address public authorizedWallet;
    mapping(address => uint) public salaries;

    event Payout(address indexed employee, uint salaryInETH);
    event PayoutFailed(address indexed employee, uint salaryInETH, string reason);

    constructor(address _authorizedWallet, address _priceFeedAddress) {
        authorizedWallet = _authorizedWallet;
        dataFeed = AggregatorV3Interface(_priceFeedAddress);
    }

    modifier onlyAuthorized() {
        require(msg.sender == authorizedWallet, 'Unauthorized');
        _;
    }

    function getUsdtSalary(address employee) public view returns (uint) {
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

    // using arrays to reduce gas
    function oneTimePayout(address[] memory employees, uint[] memory usdAmounts) external onlyAuthorized {
        require(employees.length == usdAmounts.length, "Mismatched input lengths");
        int ethToUSDT = getLatestUSDTPriceInETH();
        require(ethToUSDT > 0, 'Invalid price data');
        for (uint i = 0; i < employees.length; i++) {
            uint salaryInUSDT = usdAmounts[i];
            require(salaryInUSDT > 0, 'No salary set');
            uint salaryInETH = uint(salaryInUSDT * 1e18) / uint(ethToUSDT);
            salaryInETH = salaryInETH * 1e8;
            // Check sufficient balance
            require(
                address(this).balance >= salaryInETH,
                'Insufficient contract balance'
            );

            (bool success,) = employees[i].call{value: salaryInETH}("");
            if (success) {
                emit Payout(employees[i], salaryInETH);
            } else {
                emit PayoutFailed(employees[i], salaryInETH, "Transfer failed");
            }
        }

    }

    function setSalary(
        address employee,
        uint salaryInUSDT
    ) external onlyAuthorized {
        salaries[employee] = salaryInUSDT;
    }

    function getEmployeeSalaryInEth(address employee) public view returns (uint){
        uint salaryInUSDT = salaries[employee];
        require(salaryInUSDT > 0, 'No salary set');

        int ethToUSDT = getLatestUSDTPriceInETH();
        require(ethToUSDT > 0, 'Invalid price data');
        uint salaryInETH = uint(salaryInUSDT * 1e18) / uint(ethToUSDT);
        return salaryInETH * 1e8;
    }

    function payoutInETH(address payable employee) external onlyAuthorized {
        uint salaryInETH = getEmployeeSalaryInEth(employee);
        // Check sufficient balance
        require(
            address(this).balance >= salaryInETH,
            'Insufficient contract balance'
        );

        (bool success,) = employee.call{value: salaryInETH}("");
        if (success) {
            emit Payout(employee, salaryInETH);
        } else {
            emit PayoutFailed(employee, salaryInETH, "Transfer failed");
        }
    }

    // Fallback to receive ETH
    receive() external payable {}
}
