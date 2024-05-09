import { Injectable } from '@nestjs/common';
import { ethers } from 'ethers';
import { ConfigService } from '@nestjs/config';
import * as hre from 'hardhat';
import { BaseContractService } from './base-contract.service';

@Injectable()
export class SalariesService extends BaseContractService {
  getSalaries() {}

  async deploy() {
    const provider = await this.providerService.getProvider();

    const salary = await hre.artifacts.readArtifact('Salaries');
    const abi = salary.abi;
    const bytecode = salary.deployedBytecode;
    const signer = new ethers.Wallet(
      this.configService.getOrThrow('POLYGON_PK'),
      provider,
    );

    const salaryContract = new ethers.ContractFactory(
      abi,
      salary.bytecode,
      signer,
    );

    const myContract = await salaryContract.deploy(
      'multisig address',
      this.configService.getOrThrow('CHAINLINK_AGGREGATOR_V3'),
    );
    await myContract.waitForDeployment();

    console.log(
      '🚀 ~ HardhatService ~ deploySalaryContract ~ myContract:',
      myContract,
    );
    const address = myContract.getAddress();
    console.log('🚀 ~ SalariesService ~ deploy ~ address:', address);
  }
}
