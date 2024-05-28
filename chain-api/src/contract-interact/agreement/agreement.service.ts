import { Injectable } from '@nestjs/common';
import { BaseContractService } from '../../base/base-contract.service';
import * as hre from 'hardhat';
import { ethers } from 'ethers';
import { CHAINLINK } from '../../config/chainlink.config';
import {
  DeployAgreementDto,
  GetAgreementInfoDto,
  RequestAgreementDto,
} from './agreement.dto';
import { MultiSigWalletService } from '../multi-sig/multi-sig.service';
import { ProviderService } from '../../base/provider/provider.service';

@Injectable()
export class AgreementService extends BaseContractService {
  constructor(
    public readonly providerService: ProviderService,
    public readonly multiSigService: MultiSigWalletService,
  ) {
    super(providerService);
  }
  async deploy(dto: DeployAgreementDto, seed: string): Promise<any> {
    const { multiSigWallet } = dto;
    const { bytecode } = await hre.artifacts.readArtifact('Agreement');

    const abiCoder = ethers.AbiCoder.defaultAbiCoder();

    const abiEncodedConstructorArguments = abiCoder.encode(
      ['address', 'address', 'string', 'uint', 'address'],
      [
        CHAINLINK.AMOY.CHAINLINK_TOKEN,
        CHAINLINK.AMOY.ORACLE_ADDRESS,
        CHAINLINK.AMOY.JOB_IDS.BOOL,
        0,
        multiSigWallet,
      ],
    );
    const fullBytecode = bytecode + abiEncodedConstructorArguments.substring(2);
    const submitData = await this.multiSigService.submitTransaction(
      {
        contractAddress: multiSigWallet,
        destination: null,
        value: '0',
        data: fullBytecode,
      },
      seed,
    );
    delete submitData.data;
    return submitData;
  }

  async getResponse(dto: GetAgreementInfoDto, seed: string) {
    const { contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact('Agreement');
    const signer = await this.providerService.getSigner(seed);

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const answer = await contract.response();
    return answer.toString();
  }

  async request(dto: RequestAgreementDto, seed: string) {
    const { multiSigWallet, contractAddress, url } = dto;

    const ISubmitMultiSig = new ethers.Interface([
      'function request(string memory url)',
    ]);
    const data = ISubmitMultiSig.encodeFunctionData('request', [url]);

    return await this.multiSigService.submitTransaction(
      {
        contractAddress: multiSigWallet,
        destination: contractAddress,
        value: '0',
        data,
      },
      seed,
    );
  }
}
