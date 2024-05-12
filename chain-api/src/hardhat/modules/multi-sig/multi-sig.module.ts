import { Module } from '@nestjs/common';

import { ProviderModule } from 'src/provider/provider.module';

import { BaseContractService } from '../base-contract.service';
import { ProviderService } from 'src/provider/provider.service';
import { MultiSigWalletService } from './multi-sig.service';
import { MultiSigInteractController } from './multi-sig-interact.controller';

@Module({
  imports: [ProviderModule],
  controllers: [MultiSigInteractController],
  providers: [MultiSigWalletService],
  exports: [MultiSigWalletService],
})
export class MultiSigModule {}
