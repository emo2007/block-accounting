import { Module } from '@nestjs/common';
import { HardhatService } from './hardhat.service';
import { ProviderModule } from 'src/provider/provider.module';
import { MultiSigModule } from './multi-sig/multi-sig.module';
import { SalariesModule } from './salaries/salaries.module';

@Module({
  imports: [ProviderModule, MultiSigModule, SalariesModule],
  controllers: [],
  providers: [HardhatService],
  exports: [HardhatService, MultiSigModule, SalariesModule],
})
export class HardhatModule {}
