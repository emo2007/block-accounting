// SPDX-License-Identifier: MIT
// 0x2F9442900d067a3D37A1C2aE99462E055e32c741
pragma solidity ^0.8.7;

import {AggregatorV3Interface} from '@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol';

contract Salaries {
    AggregatorV3Interface internal dataFeed;
    address public multisigWallet;
    mapping(address => uint) public salaries;
    event Payout(address indexed employee, uint salaryInETH);
    event PayoutFailed(address indexed employee, uint salaryInETH, string reason);
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

    function payoutInETH(address payable employee) external onlyMultisig {
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

        (bool success, ) = employee.call{value: salaryInETH}("");
        if (success) {
            emit Payout(employee, salaryInETH);
        } else {
            emit PayoutFailed(employee, salaryInETH, "Transfer failed");
        }
    }

    function dummy() public pure returns (uint){
        return 1337;
    }

    // Fallback to receive ETH
    receive() external payable {}
}
