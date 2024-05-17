import { Injectable } from '@nestjs/common';
import { BaseContractService } from '../base-contract.service';
import { ethers, parseEther, TransactionReceipt } from 'ethers';
import {
  CreatePayoutDto,
  GetEmployeeSalariesDto,
  SalariesDeployDto,
  SetSalaryDto,
} from './salaries.dto';
import * as hre from 'hardhat';
import { MultiSigWalletService } from '../multi-sig/multi-sig.service';
import { ProviderService } from '../../../provider/provider.service';
import { DepositContractDto } from '../../../contract-interact/dto/multi-sig.dto';

@Injectable()
export class SalariesService extends BaseContractService {
  constructor(
    private readonly multiSigWalletService: MultiSigWalletService,
    public readonly providerService: ProviderService,
  ) {
    super(providerService);
  }
  async deploy(dto: SalariesDeployDto) {
    const { abi, bytecode } = await hre.artifacts.readArtifact('Salaries');

    const signer = await this.providerService.getSigner();

    const salaryContract = new ethers.ContractFactory(abi, bytecode, signer);

    const myContract = await salaryContract.deploy(
      dto.multiSigWallet,
      '0xF0d50568e3A7e8259E16663972b11910F89BD8e7',
    );
    await myContract.waitForDeployment();
    return await myContract.getAddress();
  }

  async getLicenseRequest() {
    const { abi } = await hre.artifacts.readArtifact(
      'LinkWellStringBytesConsumerContractExample',
    );
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(
      '0xbc3c4fed4C3A977b8868b589662270F1aEA6A777',
      abi,
      signer,
    );

    const answer: string = await contract.request();
    console.log('=>(salaries.service.ts:45) answer', answer);
    const licenseres = await this.getLicenseResponse();
    console.log('=>(salaries.service.ts:53) licenseres', licenseres);
    return answer;
  }

  async getLicenseResponse() {
    const { abi } = await hre.artifacts.readArtifact(
      'LinkWellStringBytesConsumerContractExample',
    );
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(
      '0xbc3c4fed4C3A977b8868b589662270F1aEA6A777',
      abi,
      signer,
    );

    const answer: string = await contract.responseBytes();
    console.log('=>(salaries.service.ts:45) answer', answer);
    return answer;
  }

  async getLatestUSDTPrice(contractAddress: string) {
    const { abi } = await hre.artifacts.readArtifact('Salaries');
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

    return await this.multiSigWalletService.submitTransaction({
      contractAddress: multiSigWallet,
      destination: contractAddress,
      value: '0',
      data,
    });
  }

  async getSalary(dto: GetEmployeeSalariesDto) {
    const { employeeAddress, contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact('Salaries');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const answer: BigInt = await contract.getUsdtSalary(employeeAddress);
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

    return await this.multiSigWalletService.submitTransaction({
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
