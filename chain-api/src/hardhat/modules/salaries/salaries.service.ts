import { Injectable } from '@nestjs/common';
import { BaseContractService } from '../base-contract.service';
import { ethers } from 'ethers';
import {
  GetEmployeeSalariesDto,
  SalariesDeployDto,
  SetSalaryDto,
} from './salaries.dto';
import * as hre from 'hardhat';
import { MultiSigWalletService } from '../multi-sig/multi-sig.service';
import { ProviderService } from '../../../provider/provider.service';

@Injectable()
export class SalariesService extends BaseContractService {
  constructor(
    private readonly multiSigWalletService: MultiSigWalletService,
    public readonly providerService: ProviderService,
  ) {
    super(providerService);
  }
  async deploy(dto: SalariesDeployDto): Promise<any> {
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

    const answer: string = await contract.getSalary(employeeAddress);
    return answer;
  }
}
