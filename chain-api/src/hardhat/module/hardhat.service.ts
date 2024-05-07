const hre = require('hardhat');
// import hre from 'hardhat';
import { Injectable } from '@nestjs/common';

@Injectable()
export class HardhatService {
  async deploySalaryContract() {
    // const { salaryAmount, userAddress } = req.body;

    // // Read the Solidity contract template file
    // const solidityCode = readSolidityTemplate(); // Implement this function to read the Solidity template file

    // // Replace placeholders in the Solidity contract template with provided values
    // const finalSolidityCode = replacePlaceholders(solidityCode, {
    //   salaryAmount,
    //   userAddress,
    // });

    // // Compile the Solidity contract
    // const compiledContract = await compileSolidity(finalSolidityCode);

    // // Deploy the contract
    // const deployedContract = await deployContract(compiledContract);
    const salaryC = await hre.ethers.getContractFactory('Salaries');
    const myContract = await salaryC.deploy();
    console.log(
      '🚀 ~ HardhatService ~ deploySalaryContract ~ myContract:',
      myContract,
    );
  }
}
