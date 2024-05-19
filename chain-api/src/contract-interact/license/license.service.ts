import { Injectable } from '@nestjs/common';
import * as hre from 'hardhat';
import { ethers } from 'ethers';
import { BaseContractService } from '../../base/base-contract.service';
import {
  DeployLicenseDto,
  GetLicenseInfoDto,
  GetLicenseResponseDto,
  GetShareLicense,
  RequestLicenseDto,
} from './license.dto';
import { MultiSigWalletService } from '../multi-sig/multi-sig.service';
import { ProviderService } from '../../base/provider/provider.service';
import { CHAINLINK } from '../../config/chainlink.config';

@Injectable()
export class LicenseService extends BaseContractService {
  constructor(
    private readonly multiSigService: MultiSigWalletService,
    public readonly providerService: ProviderService,
  ) {
    super(providerService);
  }
  async request(dto: RequestLicenseDto) {
    const { multiSigWallet, contractAddress } = dto;

    const ISubmitMultiSig = new ethers.Interface(['function request()']);
    const data = ISubmitMultiSig.encodeFunctionData('request');

    return await this.multiSigService.submitTransaction({
      contractAddress: multiSigWallet,
      destination: contractAddress,
      value: '0',
      data,
    });
  }

  async getTotalPayoutInUSD(dto: GetLicenseInfoDto) {
    const { contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact(
      'StreamingRightsManagement',
    );
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const answer: bigint = await contract.request();
    console.log('=>(license.service.ts:45) answer', answer);
    return answer.toString();
  }

  async deploy(dto: DeployLicenseDto) {
    console.log('=>(license.service.ts:53) dto', dto);
    const { multiSigWallet, shares, owners, payrollAddress } = dto;
    const { abi, bytecode } = await hre.artifacts.readArtifact(
      'StreamingRightsManagement',
    );
    const signer = await this.providerService.getSigner();

    const abiCoder = ethers.AbiCoder.defaultAbiCoder();

    const abiEncodedConstructorArguments = abiCoder.encode(
      [
        'address',
        'address',
        'string',
        'uint',
        'address',
        'address[]',
        'uint[]',
        'address',
      ],
      [
        CHAINLINK.AMOY.CHAINLINK_TOKEN,
        CHAINLINK.AMOY.ORACLE_ADDRESS,
        CHAINLINK.AMOY.JOB_IDS.UINT,
        0,
        multiSigWallet,
        owners,
        shares,
        payrollAddress,
      ],
    );
    const fullBytecode = bytecode + abiEncodedConstructorArguments.substring(2);
    const submitData = await this.multiSigService.submitTransaction({
      contractAddress: multiSigWallet,
      destination: null,
      value: '0',
      data: fullBytecode,
    });
    delete submitData.data;
    return submitData;
  }

  async getPayoutContract(dto: GetLicenseInfoDto) {
    const { contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact(
      'StreamingRightsManagement',
    );
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const answer: string = await contract.payoutContract();

    return answer;
  }

  async getOwners(dto: GetLicenseInfoDto) {
    const { contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact(
      'StreamingRightsManagement',
    );
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const owners: string[] = [];

    for (let i = 0; i < 10; i++) {
      try {
        const owner = await contract.owners(i);
        owners.push(owner);
      } catch (e) {
        // this.logger.error(e);
        console.log('OWNERS LIMIT');
        break;
      }
    }

    return owners;
  }

  async getShares(dto: GetShareLicense) {
    const { contractAddress, ownerAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact(
      'StreamingRightsManagement',
    );
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const answer: number = await contract.getShare(ownerAddress);
    console.log('=>(license.service.ts:135) answer', answer);

    return answer;
  }
}
