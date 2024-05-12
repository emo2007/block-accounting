import { Module } from '@nestjs/common';
import { SalariesService } from './salaries.service';
import { ProviderModule } from 'src/provider/provider.module';
import { SalariesController } from './salaries-interact.controller';
import { MultiSigModule } from '../multi-sig/multi-sig.module';

@Module({
  imports: [ProviderModule, MultiSigModule],
  controllers: [SalariesController],
  providers: [SalariesService],
  exports: [SalariesService],
})
export class SalariesModule {}
