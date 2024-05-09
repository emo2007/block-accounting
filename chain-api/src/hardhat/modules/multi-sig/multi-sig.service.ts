import { MultiSigWallet } from '../../typechain-types/contracts/MultiSigWallet';
import { Injectable } from '@nestjs/common';
import { ethers } from 'ethers';
import { ConfigService } from '@nestjs/config';
import * as hre from 'hardhat';
import { BaseContractService } from '../base-contract.service';
import { MultiSigWalletDto } from '../dto/multi-sig.dto';
import { SubmitTransactionDto } from 'src/contract-interact/dto/multi-sig.dto';

export class MultiSigWalletService extends BaseContractService {
  async deploy(dto: MultiSigWalletDto) {
    const { abi, bytecode } =
      await hre.artifacts.readArtifact('MultiSigWallet');

    const signer = await this.providerService.getSigner();

    const salaryContract = new ethers.ContractFactory(abi, bytecode, signer);

    const myContract = await salaryContract.deploy(
      dto.owners,
      dto.confirmations,
    );
    await myContract.waitForDeployment();

    console.log(
      '🚀 ~ HardhatService ~ deploySalaryContract ~ myContract:',
      myContract,
    );
    const address = myContract.getAddress();
    console.log('🚀 ~ SalariesService ~ deploy ~ address:', address);
  }

  async getOwners(address: string) {
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const multiSigContract = new ethers.Contract(address, abi);

    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(address, abi, signer);

    const owners = await contract.getOwners();

    return owners;
  }

  async submitTransaction(dto: SubmitTransactionDto) {
    const { destination, value, data, contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const multiSigContract = new ethers.Contract(contractAddress, abi);

    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);
    console.log(
      '🚀 ~ MultiSigWalletService ~ submitTransaction ~ contract:',
      contract.interface,
    );

    const tx = await contract.submitTransaction(
      destination,
      value,
      new TextEncoder().encode(data),
    );
    console.log('🚀 ~ MultiSigWalletService ~ submitTransaction ~ tx:', tx);

    return tx;
  }
}
