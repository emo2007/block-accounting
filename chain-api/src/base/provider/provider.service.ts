import { Injectable } from '@nestjs/common';
import { ethers } from 'ethers';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ProviderService {
  public provider: ethers.JsonRpcProvider;
  public networkId: number;
  private nodeUrl: string;
  constructor(private readonly configService: ConfigService) {
    this.networkId = parseInt(
      this.configService.getOrThrow('POLYGON_NETWORK_ID'),
    );
    this.nodeUrl = this.configService.getOrThrow('POLYGON_NODE');
  }

  async getProvider() {
    if (this.provider) {
      return this.provider;
    }
    this.provider = new ethers.JsonRpcProvider(this.nodeUrl, this.networkId);
    return this.provider;
  }

  async getSigner() {
    if (!this.provider) {
      await this.getProvider();
    }
    return new ethers.Wallet(
      this.configService.getOrThrow('POLYGON_PK'),
      this.provider,
    );
  }
}
