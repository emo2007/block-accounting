// const hre = require('hardhat');
import * as hre from 'hardhat';
import { Injectable } from '@nestjs/common';
import { ethers } from 'ethers';
import { ConfigService } from '@nestjs/config';
@Injectable()
export class HardhatService {
  constructor(private readonly configService: ConfigService) {}
  async deploySalaryContract() {
    const provider = new ethers.JsonRpcProvider(
      'https://polygon-amoy.g.alchemy.com/v2/pEtFFy_Qr_NrM1vMnlzSXmYXkozVNzLy',
      80002,
    );

    const salary = await hre.artifacts.readArtifact('Salaries');
    const abi = salary.abi;
    console.log('🚀 ~ HardhatService ~ deploySalaryContract ~ abi:', abi);
    const bytecode = salary.deployedBytecode;
    console.log(
      '🚀 ~ HardhatService ~ deploySalaryContract ~ bytecode:',
      bytecode,
    );
    const signer = new ethers.Wallet(
      this.configService.getOrThrow('POLYGON_PK'),
      provider,
    );

    const salaryContract = new ethers.ContractFactory(
      abi,
      salary.bytecode,
      signer,
    );

    const myContract = await salaryContract.deploy();
    await myContract.waitForDeployment();

    console.log(
      '🚀 ~ HardhatService ~ deploySalaryContract ~ myContract:',
      myContract,
    );
  }
}
