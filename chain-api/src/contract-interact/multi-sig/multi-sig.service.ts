import { ethers, parseEther, TransactionReceipt } from 'ethers';
import * as hre from 'hardhat';
import { MultiSigWalletDto } from './multi-sig.dto';
import {
  ConfirmTransactionDto,
  DepositContractDto,
  ExecuteTransactionDto,
  GetTransactionDto,
  RevokeConfirmationDto,
  SubmitTransactionDto,
} from 'src/contract-interact/multi-sig.dto';
import { parseLogs } from 'src/ethers-custom/ethers.helpers';
import { BaseContractService } from '../../base/base-contract.service';
import { getContractAddress } from '@ethersproject/address';

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
    return myContract.getAddress();
  }

  async getOwners(address: string) {
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');

    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(address, abi, signer);

    return await contract.getOwners();
  }

  async submitTransaction(dto: SubmitTransactionDto) {
    const { destination, value, data, contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const tx = await contract.submitTransaction(
      destination || '0x0000000000000000000000000000000000000000',
      value,
      data,
    );
    const txResponse: TransactionReceipt = await tx.wait();

    const eventParse = parseLogs(txResponse, contract, 'SubmitTransaction');

    return {
      txHash: txResponse.hash,
      sender: eventParse.args[0].toString(),
      txIndex: eventParse.args[1].toString(),
      to: eventParse.args[2].toString(),
      value: eventParse.args[3].toString(),
      data: eventParse.args[4].toString(),
    };
  }

  async confirmTransaction(dto: ConfirmTransactionDto) {
    const { contractAddress, index } = dto;
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const tx = await contract.confirmTransaction(index);

    const txResponse: TransactionReceipt = await tx.wait();

    const eventParse = parseLogs(txResponse, contract, 'ConfirmTransaction');

    return {
      txHash: txResponse.hash,
      sender: eventParse.args[0].toString(),
      txIndex: eventParse.args[1].toString(),
    };
  }

  async executeTransaction(dto: ExecuteTransactionDto) {
    const { index, contractAddress, isDeploy } = dto;
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const input = dto.index + new Date().getTime().toString();
    const hashed = ethers.keccak256(ethers.toUtf8Bytes(input));
    const salt = BigInt(hashed.substring(0, 10));

    if (isDeploy) {
      const tx = await contract.executeDeployTransaction(index, salt);

      const txResponse: TransactionReceipt = await tx.wait();
      const eventParse = parseLogs(txResponse, contract, 'ExecuteTransaction');
      const deployedParse = parseLogs(txResponse, contract, 'ContractDeployed');
      return {
        txHash: txResponse.hash,
        sender: eventParse.args[0].toString(),
        txIndex: eventParse.args[1].toString(),
        deployedAddress: deployedParse.args[0].toString(),
      };
    } else {
      const tx = await contract.executeTransaction(index);

      const txResponse: TransactionReceipt = await tx.wait();

      const eventParse = parseLogs(txResponse, contract, 'ExecuteTransaction');
      return {
        txHash: txResponse.hash,
        sender: eventParse.args[0].toString(),
        txIndex: eventParse.args[1].toString(),
      };
    }
  }

  async calculateFutureAddress(contractAddress: string) {
    const provider = await this.providerService.getProvider();

    const nonce = await provider.getTransactionCount(contractAddress);

    return getContractAddress({
      from: contractAddress,
      nonce: nonce,
    });
  }

  async revokeConfirmation(dto: RevokeConfirmationDto) {
    const { index, contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const tx = await contract.revokeConfirmation(index);

    return tx;
  }

  async getTransactionCount(contractAddress: string) {
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const txCount = await contract.getTransactionCount();

    return txCount;
  }

  async getTransaction(dto: GetTransactionDto) {
    const { index, contractAddress } = dto;
    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const signer = await this.providerService.getSigner();

    const contract = new ethers.Contract(contractAddress, abi, signer);

    const tx = await contract.getTransaction(index);

    return tx;
  }

  async deposit(dto: DepositContractDto) {
    const { contractAddress, value } = dto;
    const convertValue = parseEther(value);
    const signer = await this.providerService.getSigner();

    const { abi } = await hre.artifacts.readArtifact('MultiSigWallet');
    const contract = new ethers.Contract(contractAddress, abi, signer);

    const tx = await signer.sendTransaction({
      to: contractAddress,
      value: convertValue,
    });

    const txResponse: TransactionReceipt = await tx.wait();

    const eventParse = parseLogs(txResponse, contract, 'ExecuteTransaction');

    return {
      txHash: txResponse.hash,
      sender: eventParse.args[0].toString(),
      value: eventParse.args[1].toString(),
      contractBalance: eventParse.args[2].toString(),
    };
  }
}
