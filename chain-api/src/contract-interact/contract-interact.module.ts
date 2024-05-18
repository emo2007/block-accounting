import { Module } from '@nestjs/common';
import { SalariesModule } from './salaries/salaries.module';
import { MultiSigModule } from './multi-sig/multi-sig.module';
import { LicenseModule } from './license/license.module';

@Module({
  imports: [SalariesModule, MultiSigModule, LicenseModule],
  controllers: [],
  providers: [],
})
export class ContractInteractModule {}
