import { Injectable } from '@nestjs/common';
import { ethers, parseEther, TransactionReceipt } from 'ethers';
import {
  CreatePayoutDto,
  GetEmployeeSalariesDto,
  SalariesDeployDto,
  SetSalaryDto,
} from './salaries.dto';
import * as hre from 'hardhat';
import { BaseContractService } from '../../base/base-contract.service';
import { DepositContractDto } from '../multi-sig.dto';
import { CHAINLINK } from '../../config/chainlink.config';
import { ProviderService } from '../../base/provider/provider.service';
import { MultiSigWalletService } from '../multi-sig/multi-sig.service';

@Injectable()
export class SalariesService extends BaseContractService {
  constructor(
    public readonly providerService: ProviderService,
    public readonly multiSigService: MultiSigWalletService,
  ) {
    super(providerService);
  }
  async deploy(dto: SalariesDeployDto) {
    const { abi, bytecode } = await hre.artifacts.readArtifact('Payroll');

    const signer = await this.providerService.getSigner();

    const salaryContract = new ethers.ContractFactory(abi, bytecode, signer);

    const myContract = await salaryContract.deploy(
      dto.authorizedWallet,
      CHAINLINK.AMOY.AGGREGATOR_ADDRESS.USDT_ETH,
    );
    await myContract.waitForDeployment();
    return await myContract.getAddress();
  }

  async getLatestUSDTPrice(contractAddress: string) {
    const { abi } = await hre.artifacts.readArtifact('Payroll');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const answer: string = await contract.getLatestUSDTPriceInETH();
    return parseInt(answer) / 1e8;
  }

  async setSalary(dto: SetSalaryDto) {
    const { employeeAddress, salary, contractAddress, multiSigWallet } = dto;
    const ISubmitMultiSig = new ethers.Interface([
      'function setSalary(address employee, uint salaryInUSDT)',
    ]);

    const data = ISubmitMultiSig.encodeFunctionData('setSalary', [
      employeeAddress,
      salary,
    ]);

    return await this.multiSigService.submitTransaction({
      contractAddress: multiSigWallet,
      destination: contractAddress,
      value: '0',
      data,
    });
  }

  async getSalary(dto: GetEmployeeSalariesDto) {
    const { employeeAddress, contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact('Payroll');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const answer: bigint = await contract.getUsdtSalary(employeeAddress);
    return {
      salaryInUsd: answer.toString(),
    };
  }

  async createPayout(dto: CreatePayoutDto) {
    const { employeeAddress, contractAddress, multiSigWallet } = dto;
    const ISubmitMultiSig = new ethers.Interface([
      'function payoutInETH(address employee)',
    ]);
    const data = ISubmitMultiSig.encodeFunctionData('payoutInETH', [
      employeeAddress,
    ]);

    return await this.multiSigService.submitTransaction({
      contractAddress: multiSigWallet,
      destination: contractAddress,
      value: '0',
      data,
    });
  }

  async deposit(dto: DepositContractDto) {
    const { contractAddress, value } = dto;
    const signer = await this.providerService.getSigner();

    const convertValue = parseEther(value);

    const tx = await signer.sendTransaction({
      to: contractAddress,
      value: convertValue,
    });

    const txResponse: TransactionReceipt = await tx.wait();

    return txResponse;
  }
}
