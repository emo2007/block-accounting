import { Module } from '@nestjs/common';

import { ProviderModule } from 'src/provider/provider.module';

import { BaseContractService } from '../base-contract.service';
import { ProviderService } from 'src/provider/provider.service';
import { MultiSigWalletService } from './multi-sig.service';

@Module({
  imports: [ProviderModule],
  controllers: [],
  providers: [MultiSigWalletService],
  exports: [MultiSigWalletService],
})
export class MultiSigModule {}
