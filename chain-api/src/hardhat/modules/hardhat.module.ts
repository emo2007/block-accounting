import { Module } from '@nestjs/common';
import { HardhatService } from './hardhat.service';
import { ProviderModule } from 'src/provider/provider.module';
import { MultiSigWalletService } from './multi-sig/multi-sig.service';
import { SalariesService } from './salary.service';
import { BaseContractService } from './base-contract.service';
import { MultiSigModule } from './multi-sig/multi-sig.module';

@Module({
  imports: [ProviderModule, MultiSigModule],
  controllers: [],
  providers: [HardhatService, SalariesService],
  exports: [HardhatService, SalariesService, MultiSigModule],
})
export class HardhatModule {}
