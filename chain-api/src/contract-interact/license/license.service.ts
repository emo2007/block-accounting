import { Injectable } from '@nestjs/common';
import * as hre from 'hardhat';
import { ethers } from 'ethers';
import { BaseContractService } from '../../base/base-contract.service';
import {
  DeployLicenseDto,
  GetLicenseInfoDto,
  GetShareLicense,
  LicensePayoutDto,
  RequestLicenseDto,
  SetPayoutContractDto,
} from './license.dto';
import { CHAINLINK } from '../../config/chainlink.config';
import { ProviderService } from '../../base/provider/provider.service';
import { MultiSigWalletService } from '../multi-sig/multi-sig.service';

@Injectable()
export class LicenseService extends BaseContractService {
  constructor(
    public readonly providerService: ProviderService,
    public readonly multiSigService: MultiSigWalletService,
  ) {
    super(providerService);
  }
  async request(dto: RequestLicenseDto) {
    const { multiSigWallet, contractAddress, url } = dto;

    const ISubmitMultiSig = new ethers.Interface([
      'function request(string memory url)',
    ]);
    const data = ISubmitMultiSig.encodeFunctionData('request', [url]);

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

    const answer: bigint = await contract.totalPayoutInUSD();
    return answer.toString();
  }

  async deploy(dto: DeployLicenseDto) {
    const { multiSigWallet, shares, owners } = dto;
    const { bytecode } = await hre.artifacts.readArtifact(
      'StreamingRightsManagement',
    );

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
      ],
      [
        CHAINLINK.AMOY.CHAINLINK_TOKEN,
        CHAINLINK.AMOY.ORACLE_ADDRESS,
        CHAINLINK.AMOY.JOB_IDS.UINT,
        0,
        multiSigWallet,
        owners,
        shares,
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

  async payout(dto: LicensePayoutDto) {
    const { multiSigWallet, contractAddress } = dto;

    const ISubmitMultiSig = new ethers.Interface(['function payout()']);
    const data = ISubmitMultiSig.encodeFunctionData('payout');

    return await this.multiSigService.submitTransaction({
      contractAddress: multiSigWallet,
      destination: contractAddress,
      value: '0',
      data,
    });
  }

  async setPayoutContract(dto: SetPayoutContractDto) {
    const { multiSigWallet, contractAddress, payoutContract } = dto;

    const ISubmitMultiSig = new ethers.Interface([
      'function setPayoutContract(address payable)',
    ]);
    const data = ISubmitMultiSig.encodeFunctionData('setPayoutContract', [
      payoutContract,
    ]);

    return await this.multiSigService.submitTransaction({
      contractAddress: multiSigWallet,
      destination: contractAddress,
      value: '0',
      data,
    });
  }
}
