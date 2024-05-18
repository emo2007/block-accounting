import { Module } from '@nestjs/common';
import { SalariesService } from './salaries.service';
import { SalariesController } from './salaries-interact.controller';
import { MultiSigModule } from '../multi-sig/multi-sig.module';
import { BaseModule } from '../../base/base.module';

@Module({
  imports: [BaseModule, MultiSigModule],
  controllers: [SalariesController],
  providers: [SalariesService],
  exports: [SalariesService],
})
export class SalariesModule {}
