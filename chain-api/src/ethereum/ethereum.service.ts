import { Inject, Injectable } from '@nestjs/common';
import { ethers } from 'ethers';

@Injectable()
export class EthereumService {
  async getAddressFromPrivateKey(privateKey: string) {
    const wallet = new ethers.Wallet(privateKey);
    return wallet.address;
  }
  async getAddressFromSeedPhrase(seedPhrase: string) {
    const wallet = ethers.Wallet.fromPhrase(seedPhrase);
    return wallet.address;
  }
}
